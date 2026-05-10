package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"meeting-server/internal/model"
	"meeting-server/internal/repository"
)

type RecordingHandler struct {
	repo *repository.MeetingRepository
}

func NewRecordingHandler(repo *repository.MeetingRepository) *RecordingHandler {
	return &RecordingHandler{repo: repo}
}

type StartRecordingRequest struct {
	MeetingID string `json:"meeting_id"`
}

type StartRecordingResponse struct {
	RecordingID string `json:"recording_id"`
	Status      string `json:"status"`
}

func (h *RecordingHandler) StartRecording(w http.ResponseWriter, r *http.Request) {
	var req StartRecordingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	recording := &model.Recording{
		ID:        uuid.New().String(),
		MeetingID: req.MeetingID,
		Status:    "recording",
		CreatedAt: time.Now(),
	}

	json.NewEncoder(w).Encode(StartRecordingResponse{
		RecordingID: recording.ID,
		Status:      "recording",
	})
}

func (h *RecordingHandler) StopRecording(w http.ResponseWriter, r *http.Request) {
	recordingID := r.URL.Query().Get("id")
	if recordingID == "" {
		http.Error(w, "recording id required", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"recording_id": recordingID,
		"status":       "stopped",
	})
}

func (h *RecordingHandler) GetRecording(w http.ResponseWriter, r *http.Request) {
	recordingID := r.URL.Query().Get("id")
	if recordingID == "" {
		http.Error(w, "recording id required", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":     recordingID,
		"status": "ready",
	})
}