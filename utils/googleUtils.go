package utils

import (
	"backend/models"
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

func AdjustEvent(event models.CreateEvent) *calendar.Event {
	attendees := []*calendar.EventAttendee{}
	for _, a := range event.Attendees {
		attendees = append(attendees, &calendar.EventAttendee{
			Email: a.Email,
		})
	}

	return &calendar.Event{
		Summary:     event.Summary,
		Description: event.Description,
		Location:    event.Location,
		Start: &calendar.EventDateTime{
			DateTime: event.StartTime.Format(time.RFC3339),
			TimeZone: "Asia/Jakarta",
		},
		End: &calendar.EventDateTime{
			DateTime: event.EndTime.Format(time.RFC3339),
			TimeZone: "Asia/Jakarta",
		},
		Attendees: attendees,
	}
}

func ParseDateTime(dt *calendar.EventDateTime) time.Time {
	if dt.DateTime != "" {
		t, _ := time.Parse(time.RFC3339, dt.DateTime)
		return t
	}

	t, _ := time.Parse("2006-01-02", dt.Date)
	return t
}

func RefreshAccessToken(ctx context.Context, config *oauth2.Config, refreshToken string) (*oauth2.Token, error) {
	tok := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := config.TokenSource(ctx, tok)

	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	return newToken, nil
}
