package handler

import (
	"encoding/json"
	"net/http"

	"meeting-server/internal/model"
	"meeting-server/internal/service"

	"github.com/google/uuid"
)

type MeetingHandler struct {
	service *service.MeetingService
}

func NewMeetingHandler(s *service.MeetingService) *MeetingHandler {
	return &MeetingHandler{service: s}
}

func (h *MeetingHandler) CreateMeeting(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title    string `json:"title"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	meeting := &model.Meeting{
		ID:     uuid.New().String()[:8],
		Title:  req.Title,
		Status: "pending",
	}

	h.service.CreateMeeting(meeting)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meeting)
}

func (h *MeetingHandler) GetMeeting(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	meeting, err := h.service.GetMeeting(id)
	if err != nil {
		http.Error(w, "meeting not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meeting)
}

func (h *MeetingHandler) JoinMeeting(w http.ResponseWriter, r *http.Request) {
	var req struct {
		MeetingID string `json:"meeting_id"`
		UserID    string `json:"user_id"`
		Password  string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	participant := &model.Participant{
		ID:        uuid.New().String(),
		MeetingID: req.MeetingID,
		UserID:    req.UserID,
		Role:      "participant",
	}

	h.service.AddParticipant(participant)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(participant)
}

func (h *MeetingHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/meeting/create", h.CreateMeeting)
	mux.HandleFunc("/api/meeting/get", h.GetMeeting)
	mux.HandleFunc("/api/meeting/join", h.JoinMeeting)
}
