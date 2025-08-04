package adapters

import (
	"backend/models"
	"backend/utils"
	"context"
	"fmt"
	"sort"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type GoogleAdapter struct{}

func NewGoogleAdapter() *GoogleAdapter {
	return &GoogleAdapter{}
}

func (a *GoogleAdapter) newClient(ctx context.Context, accessToken string) (*calendar.Service, error) {
	token := &oauth2.Token{AccessToken: accessToken}
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	return calendar.NewService(ctx, option.WithHTTPClient(client))
}

func (a *GoogleAdapter) ListEvents(accessToken string, query models.EventQuery) ([]models.Event, error) {
	ctx := context.Background()
	srv, err := a.newClient(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	call := srv.Events.List(query.CalendarID).MaxResults(20).
		TimeMin(query.From).
		TimeMax(query.To)

	resp, err := call.Do()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var eventList []models.Event
	for _, i := range resp.Items {
		start := utils.ParseDateTime(i.Start)
		end := utils.ParseDateTime(i.End)
		eventList = append(eventList, models.Event{
			ID:          i.Id,
			Summary:     i.Summary,
			Description: i.Description,
			Location:    i.Location,
			StartTime:   start,
			EndTime:     end,
		})
	}
	sort.Slice(eventList, func(i, j int) bool {
		return eventList[i].StartTime.Before(eventList[j].StartTime)
	})

	return eventList, nil
}

func (a *GoogleAdapter) CreateEvent(accessToken string, newEvent *calendar.Event) error {
	ctx := context.Background()
	srv, err := a.newClient(ctx, accessToken)
	if err != nil {
		return err
	}

	calendarID := "primary"

	_, err = srv.Events.Insert(calendarID, newEvent).Do()
	if err != nil {
		return err
	}
	return nil
}

func (a *GoogleAdapter) UpdateEvent(accessToken string, e calendar.Event) error {
	ctx := context.Background()

	srv, err := a.newClient(ctx, accessToken)
	if err != nil {
		return err
	}

	_, err = srv.Events.Update("primary", e.Id, &e).Do()
	return err
}

func (a *GoogleAdapter) DeleteEvent(accessToken, eventID string) error {
	ctx := context.Background()
	srv, err := a.newClient(ctx, accessToken)
	if err != nil {
		return err
	}
	return srv.Events.Delete("primary", eventID).Do()
}
