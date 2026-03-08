package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/feedme/order-controller/internal/api"
	"github.com/feedme/order-controller/internal/engine"
)

func main() {
	e := engine.New(10 * time.Second)
	h := api.NewHandler(e)
	hub := api.NewHub(e)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/orders", h.CreateOrder)
	mux.HandleFunc("POST /api/bots", h.AddBot)
	mux.HandleFunc("DELETE /api/bots", h.RemoveBot)
	mux.HandleFunc("GET /api/state", h.GetState)
	mux.HandleFunc("GET /ws", hub.HandleWS)

	fmt.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
