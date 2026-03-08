package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/feedme/order-controller/internal/engine"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Hub struct {
	mu      sync.Mutex
	clients map[*websocket.Conn]bool
	engine  *engine.Engine
}

func NewHub(e *engine.Engine) *Hub {
	h := &Hub{
		clients: make(map[*websocket.Conn]bool),
		engine:  e,
	}

	e.OnEvent(func(ev engine.Event) {
		h.broadcast(ev)
	})

	return h
}

func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade error: %v", err)
		return
	}

	h.mu.Lock()
	h.clients[conn] = true
	h.mu.Unlock()

	// Send current state on connect
	state := h.engine.State()
	msg, _ := json.Marshal(engine.Event{Type: "state_sync", Payload: state})
	conn.WriteMessage(websocket.TextMessage, msg)

	// Read loop (keeps connection alive, handles close)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			h.mu.Lock()
			delete(h.clients, conn)
			h.mu.Unlock()
			conn.Close()
			break
		}
	}
}

func (h *Hub) broadcast(ev engine.Event) {
	msg, err := json.Marshal(ev)
	if err != nil {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	for conn := range h.clients {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			conn.Close()
			delete(h.clients, conn)
		}
	}
}
