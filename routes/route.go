package routes

import (
	"backend/controllers"

	"github.com/labstack/echo/v4"
)

func SetupAuthRoutes(e *echo.Echo, authController *controllers.AuthController) {
	auth := e.Group("/auth")

	auth.GET("/google/callback", authController.GoogleCallback)
}

func SetupCalenderRoutes(g *echo.Group, calenderController *controllers.CalendarController) {
	g.GET("/events", calenderController.ListEvents)
	g.POST("/events", calenderController.CreateEvent)
	g.POST("/edit/events", calenderController.UpdateEvent)
	g.POST("/delete/events/:id", calenderController.DeleteEvent)
}
