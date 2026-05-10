package repository

import (
    "github.com/jmoiron/sqlx"
    "meeting-server/internal/model"
)

type MeetingRepository struct {
    db *sqlx.DB
}

func NewMeetingRepository(db *sqlx.DB) *MeetingRepository {
    return &MeetingRepository{db: db}
}

func (r *MeetingRepository) Create(meeting *model.Meeting) error {
    query := `INSERT INTO meetings (id, title, password_hash, host_id, status) VALUES ($1, $2, $3, $4, $5)`
    _, err := r.db.Exec(query, meeting.ID, meeting.Title, meeting.Password, meeting.HostID, meeting.Status)
    return err
}

func (r *MeetingRepository) GetByID(id string) (*model.Meeting, error) {
    var meeting model.Meeting
    query := `SELECT * FROM meetings WHERE id = $1`
    err := r.db.Get(&meeting, query, id)
    if err != nil {
        return nil, err
    }
    return &meeting, nil
}

func (r *MeetingRepository) UpdateStatus(id, status string) error {
    query := `UPDATE meetings SET status = $1 WHERE id = $2`
    _, err := r.db.Exec(query, status, id)
    return err
}

func (r *MeetingRepository) AddParticipant(p *model.Participant) error {
    query := `INSERT INTO meeting_participants (id, meeting_id, user_id, role) VALUES ($1, $2, $3, $4)`
    _, err := r.db.Exec(query, p.ID, p.MeetingID, p.UserID, p.Role)
    return err
}

func (r *MeetingRepository) GetParticipants(meetingID string) ([]*model.Participant, error) {
    var participants []*model.Participant
    query := `SELECT * FROM meeting_participants WHERE meeting_id = $1`
    err := r.db.Select(&participants, query, meetingID)
    return participants, err
}

func (r *MeetingRepository) RemoveParticipant(meetingID, userID string) error {
    query := `DELETE FROM meeting_participants WHERE meeting_id = $1 AND user_id = $2`
    _, err := r.db.Exec(query, meetingID, userID)
    return err
}