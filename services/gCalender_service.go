package services

import (
	"backend/adapters"
	"backend/models"
	"backend/utils"
	"fmt"
	"time"
)

type CalendarService struct {
	adapter *adapters.GoogleAdapter
}

func NewCalendarService(adapter *adapters.GoogleAdapter) *CalendarService {
	return &CalendarService{adapter: adapter}
}

func (s *CalendarService) ListEvents(token string, query models.EventQuery) ([]models.Event, error) {
	now := time.Now()
	rfc3339Time := now.Format(time.RFC3339)

	tenYearsLater := now.AddDate(10, 0, 0)
	tenYearsLaterString := tenYearsLater.Format(time.RFC3339)
	if query.From == "" {
		query.From = rfc3339Time
	}

	if query.To == "" {
		query.To = tenYearsLaterString
	}
	return s.adapter.ListEvents(token, query)
}

func (s *CalendarService) Create(token string, newEvents []models.CreateEvent) error {
	var failedEvents []string
	for _, e := range newEvents {
		eventToInsert := utils.AdjustEvent(e)
		err := s.adapter.CreateEvent(token, eventToInsert)
		if err != nil {
			failedEvents = append(failedEvents, e.Summary)
			continue
		}
	}
	if len(failedEvents) > 0 {
		return fmt.Errorf("Failed to create the following events: %v", failedEvents)
	}

	return nil
}

func (s *CalendarService) Update(token string, e models.EditEvent) error {
	if e.ID == "" {
		return fmt.Errorf("event ID is required")
	}

	if e.StartTime.IsZero() || e.EndTime.IsZero() {
		return fmt.Errorf("start time and end time are required")
	}

	event := models.CreateEvent{
		Summary:     e.Summary,
		Description: e.Description,
		Location:    e.Location,
		StartTime:   e.StartTime,
		EndTime:     e.EndTime,
		Attendees:   e.Attendees,
	}

	ev := utils.AdjustEvent(event)
	ev.Id = e.ID

	return s.adapter.UpdateEvent(token, *ev)
}

func (s *CalendarService) Delete(token, eventID string) error {
	return s.adapter.DeleteEvent(token, eventID)
}
