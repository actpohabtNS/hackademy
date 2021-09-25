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

func getCakeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("cake"))
	if err != nil {
		return
	}
}

func main() {
	r := mux.NewRouter()

	userService := UserService{
		repository: NewInMemoryUserStorage(),
	}

	r.HandleFunc("/cake", getCakeHandler).Methods(http.MethodGet)
	r.HandleFunc("/user/register", userService.Register).Methods(http.MethodPost)

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
