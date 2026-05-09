package model

import "time"

type Meeting struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"` // pending, active, ended
}

type Participant struct {
	ID        string `json:"id"`
	MeetingID string `json:"meeting_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"` // host, participant
}
