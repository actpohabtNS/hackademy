package main

import (
	"log"
	"net/http"
)

func getCakeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("cake"))
}
func main() {
	http.HandleFunc("/cake", getCakeHandler)
	log.Println("Server started, hit Ctrl+C to stop")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Println("Server exited with error:", err)
	}
	log.Println("Good bye :)")
}
