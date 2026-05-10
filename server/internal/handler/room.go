package handler

import (
	"encoding/json"
	"net/http"

	"meeting-server/internal/sfu"
)

type RoomHandler struct {
	sfuClient *sfu.Client
}

func NewRoomHandler(sfuClient *sfu.Client) *RoomHandler {
	return &RoomHandler{sfuClient: sfuClient}
}

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MeetingID string `json:"meeting_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	room, err := h.sfuClient.CreateRoom(req.MeetingID)
	if err != nil {
		http.Error(w, "failed to create room", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(room)
}

func (h *RoomHandler) GetRoom(w http.ResponseWriter, r *http.Request) {
	// 获取房间信息
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}