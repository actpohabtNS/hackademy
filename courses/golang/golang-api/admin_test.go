package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestAdmin_JWT(t *testing.T) {
	doRequest := createRequester(t)
	t.Run("getting supadmin jwt", func(t *testing.T) {
		u := newTestUserService()

		jwtService, jwtErr := NewJWTService("pubkey.rsa", "privkey.rsa")
		if jwtErr != nil {
			panic(jwtErr)
		}

		ts := httptest.NewServer(newRouter(u, jwtService))
		defer ts.Close()

		jwtParams := map[string]interface{}{
			"email":    os.Getenv("CAKE_SUPERADMIN_EMAIL"),
			"password": os.Getenv("CAKE_SUPERADMIN_PASSWORD"),
		}
		resp := doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/jwt", prepareParams(t, jwtParams)))
		assertStatus(t, 200, resp)
	})

	t.Run("deny regular user access admin api", func(t *testing.T) {
		u := newTestUserService()

		jwtService, jwtErr := NewJWTService("pubkey.rsa", "privkey.rsa")
		if jwtErr != nil {
			panic(jwtErr)
		}

		ts := httptest.NewServer(newRouter(u, jwtService))
		defer ts.Close()

		// registration
		registerParams := map[string]interface{}{
			"email":         "test@mail.com",
			"password":      "somepass",
			"favorite_cake": "cheesecake",
		}
		doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/register", prepareParams(t, registerParams)))

		// JWT generation
		jwtParams := map[string]interface{}{
			"email":    "test@mail.com",
			"password": "somepass",
		}
		jwtResp := doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/jwt", prepareParams(t, jwtParams)))

		// trying ban user
		banParams := map[string]interface{}{
			"email":  "whatever",
			"reason": "whatever",
		}
		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/admin/ban", prepareParams(t, banParams))
		req.Header.Set("Authorization", "Bearer "+string(jwtResp.body))
		resp := doRequest(req, nil)
		assertStatus(t, 401, resp)
		assertBody(t, "permission denied", resp)
	})

	t.Run("banning user", func(t *testing.T) {
		u := newTestUserService()

		jwtService, jwtErr := NewJWTService("pubkey.rsa", "privkey.rsa")
		if jwtErr != nil {
			panic(jwtErr)
		}

		ts := httptest.NewServer(newRouter(u, jwtService))
		defer ts.Close()

		// registration
		registerParams := map[string]interface{}{
			"email":         "test@mail.com",
			"password":      "somepass",
			"favorite_cake": "cheesecake",
		}
		doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/register", prepareParams(t, registerParams)))

		// superadmin JWT generation
		jwtParams := map[string]interface{}{
			"email":    os.Getenv("CAKE_SUPERADMIN_EMAIL"),
			"password": os.Getenv("CAKE_SUPERADMIN_PASSWORD"),
		}
		jwtResp := doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/jwt", prepareParams(t, jwtParams)))

		// trying ban user
		banParams := map[string]interface{}{
			"email":  "test@mail.com",
			"reason": "bad boy",
		}
		req, _ := http.NewRequest(http.MethodPost, ts.URL+"/admin/ban", prepareParams(t, banParams))
		req.Header.Set("Authorization", "Bearer "+string(jwtResp.body))
		resp := doRequest(req, nil)

		assertStatus(t, 200, resp)
		assertBody(t, "user test@mail.com banned", resp)

		// banned user JWT generation
		bannedJwtParams := map[string]interface{}{
			"email":    "test@mail.com",
			"password": "somepass",
		}
		bannedJwtResp := doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/jwt", prepareParams(t, bannedJwtParams)))

		assertStatus(t, 401, bannedJwtResp)
		assertBody(t, "you are banned! Reason: bad boy", bannedJwtResp)
	})

	t.Run("unbanning user", func(t *testing.T) {
		u := newTestUserService()

		jwtService, jwtErr := NewJWTService("pubkey.rsa", "privkey.rsa")
		if jwtErr != nil {
			panic(jwtErr)
		}

		ts := httptest.NewServer(newRouter(u, jwtService))
		defer ts.Close()

		// registration
		registerParams := map[string]interface{}{
			"email":         "test@mail.com",
			"password":      "somepass",
			"favorite_cake": "cheesecake",
		}
		doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/register", prepareParams(t, registerParams)))

		// superadmin JWT generation
		jwtParams := map[string]interface{}{
			"email":    os.Getenv("CAKE_SUPERADMIN_EMAIL"),
			"password": os.Getenv("CAKE_SUPERADMIN_PASSWORD"),
		}
		jwtResp := doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/jwt", prepareParams(t, jwtParams)))

		// trying ban user
		banParams := map[string]interface{}{
			"email":  "test@mail.com",
			"reason": "bad boy",
		}
		banReq, _ := http.NewRequest(http.MethodPost, ts.URL+"/admin/ban", prepareParams(t, banParams))
		banReq.Header.Set("Authorization", "Bearer "+string(jwtResp.body))
		doRequest(banReq, nil)

		// trying unban user
		unbanParams := map[string]interface{}{
			"email": "test@mail.com",
		}
		unbanReq, _ := http.NewRequest(http.MethodPost, ts.URL+"/admin/unban", prepareParams(t, unbanParams))
		unbanReq.Header.Set("Authorization", "Bearer "+string(jwtResp.body))
		doRequest(unbanReq, nil)

		// unbanned user JWT generation
		unbannedJwtParams := map[string]interface{}{
			"email":    "test@mail.com",
			"password": "somepass",
		}
		unbannedJwtResp := doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/jwt", prepareParams(t, unbannedJwtParams)))

		assertStatus(t, 200, unbannedJwtResp)
	})

	t.Run("inspecting user with ban history", func(t *testing.T) {
		u := newTestUserService()

		jwtService, jwtErr := NewJWTService("pubkey.rsa", "privkey.rsa")
		if jwtErr != nil {
			panic(jwtErr)
		}

		ts := httptest.NewServer(newRouter(u, jwtService))
		defer ts.Close()

		// registration
		registerParams := map[string]interface{}{
			"email":         "test@mail.com",
			"password":      "somepass",
			"favorite_cake": "cheesecake",
		}
		doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/register", prepareParams(t, registerParams)))

		// superadmin JWT generation
		jwtParams := map[string]interface{}{
			"email":    os.Getenv("CAKE_SUPERADMIN_EMAIL"),
			"password": os.Getenv("CAKE_SUPERADMIN_PASSWORD"),
		}
		jwtResp := doRequest(http.NewRequest(http.MethodPost, ts.URL+"/user/jwt", prepareParams(t, jwtParams)))

		// trying ban user
		banParams := map[string]interface{}{
			"email":  "test@mail.com",
			"reason": "bad boy",
		}
		banReq, _ := http.NewRequest(http.MethodPost, ts.URL+"/admin/ban", prepareParams(t, banParams))
		banReq.Header.Set("Authorization", "Bearer "+string(jwtResp.body))
		banTime := time.Now().Format(banTimeFormat)
		doRequest(banReq, nil)
		banStr := "-- was banned (reason: bad boy) at " + banTime + " by " +
			os.Getenv("CAKE_SUPERADMIN_EMAIL") + "\n"

		// trying unban user
		unbanParams := map[string]interface{}{
			"email": "test@mail.com",
		}
		unbanReq, _ := http.NewRequest(http.MethodPost, ts.URL+"/admin/unban", prepareParams(t, unbanParams))
		unbanReq.Header.Set("Authorization", "Bearer "+string(jwtResp.body))
		unbanTime := time.Now().Format(banTimeFormat)
		doRequest(unbanReq, nil)
		unbanStr := "-- was unbanned at " + unbanTime + " by " +
			os.Getenv("CAKE_SUPERADMIN_EMAIL") + "\n"

		// inspecting user
		inspectReq, _ := http.NewRequest(http.MethodGet,
			ts.URL+"/admin/inspect?email=test@mail.com",
			prepareParams(t, jwtParams))
		inspectReq.Header.Set("Authorization", "Bearer "+string(jwtResp.body))
		inspectRest := doRequest(inspectReq, nil)

		assertStatus(t, 200, inspectRest)
		assertBody(t, "user test@mail.com:\n"+banStr+unbanStr, inspectRest)
	})
}
