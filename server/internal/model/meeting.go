package model

import "time"

type Meeting struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Password  string    `json:"-"`
	HostID    string    `json:"host_id"`
	Status    string    `json:"status"` // pending, active, ended
	CreatedAt time.Time `json:"created_at"`
}

type Participant struct {
	ID        string `json:"id"`
	MeetingID string `json:"meeting_id"`
	UserID    string `json:"user_id"`
	Role      string `json:"role"` // host, participant
}
