package repository

import (
	"sync"

	"meeting-server/internal/model"
)

type MeetingRepository struct {
	meetings     map[string]*model.Meeting
	participants map[string][]*model.Participant
	mu           sync.RWMutex
}

func NewMeetingRepository() *MeetingRepository {
	return &MeetingRepository{
		meetings:     make(map[string]*model.Meeting),
		participants: make(map[string][]*model.Participant),
	}
}

func (r *MeetingRepository) Create(meeting *model.Meeting) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.meetings[meeting.ID] = meeting
}

func (r *MeetingRepository) GetByID(id string) (*model.Meeting, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if meeting, ok := r.meetings[id]; ok {
		return meeting, nil
	}
	return nil, nil
}

func (r *MeetingRepository) AddParticipant(p *model.Participant) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.participants[p.MeetingID] = append(r.participants[p.MeetingID], p)
}

func (r *MeetingRepository) GetParticipants(meetingID string) []*model.Participant {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.participants[meetingID]
}
