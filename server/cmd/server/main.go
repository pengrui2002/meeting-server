package main

import (
	"log"
	"net/http"

	"meeting-server/internal/signal"
)

func main() {
	hub := signal.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", hub.HandleWS)

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
