package signal

import (
	"log"
	"net/http"

	"meeting-server/internal/sfu"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type BroadcastMessage struct {
	RoomID string
	Data   []byte
}

type Hub struct {
	rooms      map[string]*Room
	register   chan *Client
	unregister chan *Client
	broadcast  chan BroadcastMessage
	sfuClient  *sfu.Client
}

func NewHub(sfuClient *sfu.Client) *Hub {
	return &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan BroadcastMessage),
		sfuClient:  sfuClient,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			room, ok := h.rooms[client.roomID]
			if !ok {
				room = NewRoom(client.roomID, h.sfuClient)
				h.rooms[client.roomID] = room
				go room.Run()
			}
			room.register <- client

		case client := <-h.unregister:
			if room, ok := h.rooms[client.roomID]; ok {
				room.unregister <- client
			}

		case msg := <-h.broadcast:
			if room, ok := h.rooms[msg.RoomID]; ok {
				room.broadcast <- msg.Data
			}
		}
	}
}

func (h *Hub) HandleWS(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room")
	userID := r.URL.Query().Get("user")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade error: %v", err)
		return
	}

	client := &Client{
		hub:    h,
		conn:   conn,
		send:   make(chan []byte, 256),
		roomID: roomID,
		userID: userID,
	}

	h.register <- client

	go client.WritePump()
	go client.ReadPump()
}
