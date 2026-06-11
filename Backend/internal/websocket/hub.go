package websocket

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Hub struct {
	// rooms maps showtimeID to a list of client connections
	rooms map[string]map[*websocket.Conn]bool
	mu    sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]map[*websocket.Conn]bool),
	}
}

type SeatUpdate struct {
	ShowtimeID string   `json:"showtime_id"`
	SeatIDs    []string `json:"seat_ids"`
	Status     string   `json:"status"`
	LockedBy   string   `json:"locked_by,omitempty"`
}

func (h *Hub) AddClient(showtimeID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[showtimeID] == nil {
		h.rooms[showtimeID] = make(map[*websocket.Conn]bool)
	}
	h.rooms[showtimeID][conn] = true
	log.Printf("Client connected to room: %s", showtimeID)
}

func (h *Hub) RemoveClient(showtimeID string, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[showtimeID] != nil {
		delete(h.rooms[showtimeID], conn)
		if len(h.rooms[showtimeID]) == 0 {
			delete(h.rooms, showtimeID)
		}
	}
	conn.Close()
	log.Printf("Client disconnected from room: %s", showtimeID)
}

func (h *Hub) Broadcast(update SeatUpdate) {
	h.mu.RLock()
	clients := h.rooms[update.ShowtimeID]
	h.mu.RUnlock()

	if clients == nil {
		return
	}

	message, err := json.Marshal(update)
	if err != nil {
		log.Println("Error marshalling broadcast:", err)
		return
	}

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Printf("Error sending message to client: %v", err)
			go h.RemoveClient(update.ShowtimeID, client)
		}
	}
}
