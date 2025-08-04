package models

import "time"

type Event struct {
	ID          string    `json:"id,omitempty"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

type EventQuery struct {
	CalendarID string `json:"calendar_id"`
	From       string `json:"from,omitempty"`
	To         string `json:"to,omitempty"`
}

type CreateEvent struct {
	Summary     string      `json:"summary"`
	Description string      `json:"description,omitempty"`
	Location    string      `json:"location,omitempty"`
	StartTime   time.Time   `json:"start_time"`
	EndTime     time.Time   `json:"end_time"`
	Attendees   []Attendees `json:"attendees,omitempty"`
}

type EditEvent struct {
	ID          string      `json:"id,omitempty"`
	Summary     string      `json:"summary,omitempty"`
	Description string      `json:"description,omitempty"`
	Location    string      `json:"location,omitempty"`
	StartTime   time.Time   `json:"start_time,omitempty"`
	EndTime     time.Time   `json:"end_time,omitempty"`
	Attendees   []Attendees `json:"attendees,omitempty"`
}

type Attendees struct {
	Email string `json:"email"`
}
