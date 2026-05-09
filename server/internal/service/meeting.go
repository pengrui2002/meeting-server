package service

import (
	"meeting-server/internal/model"
	"meeting-server/internal/repository"
)

type MeetingService struct {
	repo *repository.MeetingRepository
}

func NewMeetingService(r *repository.MeetingRepository) *MeetingService {
	return &MeetingService{repo: r}
}

func (s *MeetingService) CreateMeeting(meeting *model.Meeting) {
	s.repo.Create(meeting)
}

func (s *MeetingService) GetMeeting(id string) (*model.Meeting, error) {
	return s.repo.GetByID(id)
}

func (s *MeetingService) AddParticipant(p *model.Participant) {
	s.repo.AddParticipant(p)
}
