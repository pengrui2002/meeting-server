package model

import "time"

type User struct {
    ID           string    `json:"id" db:"id"`
    Username     string    `json:"username" db:"username"`
    PasswordHash string    `json:"-" db:"password_hash"`
    DisplayName  string    `json:"display_name" db:"display_name"`
    CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Recording struct {
    ID              string    `json:"id" db:"id"`
    MeetingID       string    `json:"meeting_id" db:"meeting_id"`
    HostID          string    `json:"host_id" db:"host_id"`
    Filename        string    `json:"filename" db:"filename"`
    Filepath        string    `json:"filepath" db:"filepath"`
    DurationSeconds int       `json:"duration_seconds" db:"duration_seconds"`
    FileSizeBytes   int64     `json:"file_size_bytes" db:"file_size_bytes"`
    Status          string    `json:"status" db:"status"`
    CreatedAt       time.Time `json:"created_at" db:"created_at"`
}