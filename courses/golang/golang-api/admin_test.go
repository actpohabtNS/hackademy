package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
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
}
