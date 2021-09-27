package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func wrapJwt(jwt *JWTService, f func(http.ResponseWriter, *http.Request, *JWTService)) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		f(rw, r, jwt)
	}
}

func getCakeHandler(w http.ResponseWriter, r *http.Request, u User) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("[" + u.Email + "], your favourite cake is " + u.FavoriteCake))
}

func newRouter(u *UserService, jwtService *JWTService) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/user/register", u.Register).Methods(http.MethodPost)
	r.HandleFunc("/user/me", jwtService.jwtAuth(u.repository, getCakeHandler)).Methods(http.MethodGet)
	r.HandleFunc("/user/jwt", wrapJwt(jwtService, u.JWT)).Methods(http.MethodPost)

	return r
}

func main() {
	users := NewInMemoryUserStorage()
	userService := UserService{repository: users}

	jwtService, jwtErr := NewJWTService("pubkey.rsa", "privkey.rsa")
	if jwtErr != nil {
		panic(jwtErr)
	}

	r := newRouter(&userService, jwtService)

	srv := http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		ctx, cancel := context.WithTimeout(context.Background(),
			5*time.Second)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			return
		}
	}()
	log.Println("Server started, hit Ctrl+C to stop")
	err := srv.ListenAndServe()
	if err != nil {
		log.Println("Server exited with error:", err)
	}
	log.Println("Good bye :)")
}
