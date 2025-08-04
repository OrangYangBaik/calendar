package utils

import (
	"strings"

	"github.com/labstack/echo/v4"
)

func GetBearerToken(c echo.Context) string {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check if it starts with "Bearer "
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}

	// Extract token (remove "Bearer " prefix)
	return strings.TrimPrefix(authHeader, "Bearer ")
}
