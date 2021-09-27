package main

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type parsedResponse struct {
	status int
	body   []byte
}

func createRequester(t *testing.T) func(req *http.Request, err error) parsedResponse {
	return func(req *http.Request, err error) parsedResponse {
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			return parsedResponse{}
		}

		res, httpErr := http.DefaultClient.Do(req)
		if httpErr != nil {
			t.Errorf("unexpected error: %v", httpErr)
			return parsedResponse{}
		}

		resp, ioErr := io.ReadAll(res.Body)
		closeErr := res.Body.Close()
		if closeErr != nil {
			return parsedResponse{}
		}
		if ioErr != nil {
			t.Errorf("unexpected error: %v", ioErr)
			return parsedResponse{}
		}

		return parsedResponse{res.StatusCode, resp}
	}
}

func prepareParams(t *testing.T, params map[string]interface{}) io.Reader {
	body, err := json.Marshal(params)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	return bytes.NewBuffer(body)
}

func newTestUserService() *UserService {
	return &UserService{
		repository: NewInMemoryUserStorage(),
	}
}

func newTestRouter(u *UserService, jwtService *JWTService) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/user/register", u.Register).Methods(http.MethodPost)
	r.HandleFunc("/cake", jwtService.jwtAuth(u.repository, getCakeHandler)).Methods(http.MethodGet)
	r.HandleFunc("/user/jwt", wrapJwt(jwtService, u.JWT)).Methods(http.MethodPost)

	return r
}

func assertStatus(t *testing.T, expected int, r parsedResponse) {
	if r.status != expected {
		t.Errorf("Unexpected response status. Expected: %d,actual: %d", expected, r.status)
	}
}

func assertBody(t *testing.T, expected string, r parsedResponse) {
	actual := string(r.body)
	if actual != expected {
		t.Errorf("Unexpected response body. Expected: %s,actual: %s", expected, actual)
	}
}

func TestUsers_JWT(t *testing.T) {
	doRequest := createRequester(t)
	t.Run("user does not exist", func(t *testing.T) {
		u := newTestUserService()
		j, err := NewJWTService("pubkey.rsa", "privkey.rsa")
		if err != nil {
			t.FailNow()
		}
		ts := httptest.NewServer(wrapJwt(j, u.JWT))
		defer ts.Close()
		params := map[string]interface{}{
			"email":    "test@mail.com",
			"password": "somepass",
		}
		resp := doRequest(http.NewRequest(http.MethodPost, ts.URL, prepareParams(t, params)))
		assertStatus(t, 422, resp)
		assertBody(t, "invalid login credentials", resp)
	})

	t.Run("registration", func(t *testing.T) {
		u := newTestUserService()
		ts := httptest.NewServer(logRequest(u.Register))
		defer ts.Close()
		params := map[string]interface{}{
			"email":         "test@mail.com",
			"password":      "somepass",
			"favorite_cake": "cheesecake",
		}
		resp := doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/register", prepareParams(t, params)))
		assertStatus(t, 201, resp)
		assertBody(t, "registered", resp)
	})

	t.Run("wrong password", func(t *testing.T) {
		u := newTestUserService()

		jwtService, jwtErr := NewJWTService("pubkey.rsa", "privkey.rsa")
		if jwtErr != nil {
			panic(jwtErr)
		}

		ts := httptest.NewServer(newTestRouter(u, jwtService))
		defer ts.Close()

		registerParams := map[string]interface{}{
			"email":         "test@mail.com",
			"password":      "somepass",
			"favorite_cake": "cheesecake",
		}
		doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/register", prepareParams(t, registerParams)))

		jwtParams := map[string]interface{}{
			"email":    "test@mail.com",
			"password": "wrongpass",
		}
		resp := doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/jwt", prepareParams(t, jwtParams)))
		assertStatus(t, 422, resp)
		assertBody(t, "invalid login credentials", resp)
	})

	t.Run("unauthorized", func(t *testing.T) {
		u := newTestUserService()

		jwtService, jwtErr := NewJWTService("pubkey.rsa", "privkey.rsa")
		if jwtErr != nil {
			panic(jwtErr)
		}

		ts := httptest.NewServer(newTestRouter(u, jwtService))
		defer ts.Close()

		registerParams := map[string]interface{}{
			"email":         "test@mail.com",
			"password":      "somepass",
			"favorite_cake": "cheesecake",
		}
		doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/register", prepareParams(t, registerParams)))

		resp := doRequest(http.NewRequest(http.MethodGet, ts.URL+"/cake", nil))
		assertStatus(t, 401, resp)
		assertBody(t, "unauthorized", resp)
	})

	t.Run("get cake", func(t *testing.T) {
		u := newTestUserService()

		jwtService, jwtErr := NewJWTService("pubkey.rsa", "privkey.rsa")
		if jwtErr != nil {
			panic(jwtErr)
		}

		ts := httptest.NewServer(newTestRouter(u, jwtService))
		defer ts.Close()

		registerParams := map[string]interface{}{
			"email":         "test@mail.com",
			"password":      "somepass",
			"favorite_cake": "cheesecake",
		}
		doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/register", prepareParams(t, registerParams)))

		jwtParams := map[string]interface{}{
			"email":    "test@mail.com",
			"password": "somepass",
		}
		jwtResp := doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/jwt", prepareParams(t, jwtParams)))
		req, _ := http.NewRequest(http.MethodGet, ts.URL+"/cake", nil)
		req.Header.Set("Authorization", "Bearer "+string(jwtResp.body))
		resp := doRequest(req, nil)
		assertStatus(t, 200, resp)
		assertBody(t, "cheesecake", resp)
	})
}
