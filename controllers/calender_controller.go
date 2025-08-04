package controllers

import (
	"backend/models"
	"backend/services"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CalendarController struct {
	Svc *services.CalendarService
}

func NewCalendarController(svc *services.CalendarService) *CalendarController {
	return &CalendarController{Svc: svc}
}

func (c *CalendarController) ListEvents(ctx echo.Context) error {
	var eventParam models.EventQuery
	// token := utils.GetBearerToken(ctx)
	// if token == "" {
	// 	return ctx.JSON(http.StatusUnauthorized, echo.Map{
	// 		"error": "Authorization header with Bearer token required",
	// 	})
	// }

	accessToken := ctx.Get("googleAccessToken").(string)
	if accessToken == "" {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Authorization header with Bearer token required",
		})
	}

	if err := ctx.Bind(&eventParam); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid query parameters",
		})
	}

	events, err := c.Svc.ListEvents(accessToken, eventParam)
	if err != nil {
		fmt.Println("\nerror:", err)
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, events)
}

func (c *CalendarController) CreateEvent(ctx echo.Context) error {
	// token := utils.GetBearerToken(ctx)
	// if token == "" {
	// 	return ctx.JSON(http.StatusUnauthorized, echo.Map{
	// 		"error": "Authorization header with Bearer token required",
	// 	})
	// }

	accessToken := ctx.Get("googleAccessToken").(string)
	if accessToken == "" {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Authorization header with Bearer token required",
		})
	}

	var newEvents []models.CreateEvent
	if err := ctx.Bind(&newEvents); err != nil {
		fmt.Println(err.Error())
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	err := c.Svc.Create(accessToken, newEvents)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, "all events successfully created")
}

func (c *CalendarController) UpdateEvent(ctx echo.Context) error {
	var eventParam models.EditEvent

	accessToken := ctx.Get("googleAccessToken").(string)
	if accessToken == "" {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Authorization header with Bearer token required",
		})
	}

	if err := ctx.Bind(&eventParam); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{
			"error":         "invalid query parameters",
			"error message": err.Error(),
		})
	}

	if err := c.Svc.Update(accessToken, eventParam); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, "all events successfully edited")
}

func (c *CalendarController) DeleteEvent(ctx echo.Context) error {
	accessToken := ctx.Get("googleAccessToken").(string)
	if accessToken == "" {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Authorization header with Bearer token required",
		})
	}

	eventID := ctx.Param("id")
	if eventID == "" {
		return ctx.JSON(http.StatusBadRequest, "accessToken and event id required")
	}

	fmt.Println("eventId: ", eventID)
	if err := c.Svc.Delete(accessToken, eventID); err != nil {
		return ctx.JSON(http.StatusInternalServerError, err.Error())
	}
	return ctx.JSON(http.StatusOK, "Event successfully deleted")
}
