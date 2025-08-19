package controllers

import (
	"backend/constants"
	"backend/dtos"
	"backend/services"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthController struct {
	GoogleOAuthConfig *oauth2.Config
	authService       services.AuthService
}

func NewAuthController(authService services.AuthService) *AuthController {
	return &AuthController{
		GoogleOAuthConfig: &oauth2.Config{
			ClientID:     os.Getenv(constants.GOOGLE_CLIENT_ID),
			ClientSecret: os.Getenv(constants.GOOGLE_CLIENT_SECRET),
			//RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			RedirectURL: "postmessage",
			Scopes: []string{
				"https://www.googleapis.com/auth/calendar",
			},
			Endpoint: google.Endpoint,
		},
		authService: authService,
	}
}

func (ac *AuthController) GoogleCallback(c echo.Context) error {
	var jwtToken string
	var expiresAt time.Time
	// Verify state (CSRF protection)
	// state := c.QueryParam("state")
	// stateCookie, err := c.Cookie(constants.StateSessionKey)
	// if err != nil || stateCookie.Value != state {
	// 	return c.JSON(http.StatusBadRequest, map[string]string{
	// 		"error": "Invalid state token",
	// 	})
	// }

	code := c.QueryParam("code")
	if code == "" {
		errorMsg := fmt.Errorf("Authorization code not found")
		return c.JSON(http.StatusBadRequest, dtos.Response{
			Data:  nil,
			Error: errorMsg.Error(),
		})
	}

	token, err := ac.GoogleOAuthConfig.Exchange(
		c.Request().Context(),
		code,
		oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("prompt", "consent"),
		oauth2.SetAuthURLParam("redirect_uri", ac.GoogleOAuthConfig.RedirectURL),
	)
	if err != nil {
		errorMsg := fmt.Errorf("Failed to exchange token: %v\n", err)
		return c.JSON(http.StatusBadRequest, dtos.Response{
			Data:  nil,
			Error: errorMsg.Error(),
		})
	}

	userInfo, err := ac.authService.GetUserInfo(token.AccessToken)
	if err != nil {
		errorMsg := fmt.Errorf("Error getting user info: %v\n", err)
		return c.JSON(http.StatusBadRequest, dtos.Response{
			Data:  nil,
			Error: errorMsg.Error(),
		})
	}

	user, err := ac.authService.ProcessGoogleUser(userInfo, token, "")
	if err != nil {
		errorMsg := fmt.Errorf("Error processing user: %v\n", err)
		return c.JSON(http.StatusBadRequest, dtos.Response{
			Data:  nil,
			Error: errorMsg.Error(),
		})
	}

	jwtToken, expiresAt, err = ac.authService.GenerateJWT(user)
	if err != nil {
		errorMsg := fmt.Errorf("Error generating JWT: %v\n", err)
		return c.JSON(http.StatusBadRequest, dtos.Response{
			Data:  nil,
			Error: errorMsg.Error(),
		})
	}

	if user.FolderID == "" {
		newUser, err := ac.authService.ProcessGoogleUser(userInfo, token, jwtToken)
		if err != nil {
			errorMsg := fmt.Errorf("Failed to create new user: %v\n", err)
			return c.JSON(http.StatusInternalServerError, dtos.Response{
				Data:  nil,
				Error: errorMsg.Error(),
			})
		}

		jwtToken, expiresAt, err = ac.authService.GenerateJWT(newUser)
		if err != nil {
			errorMsg := fmt.Errorf("Error generating JWT: %v\n", err)
			return c.JSON(http.StatusBadRequest, dtos.Response{
				Data:  nil,
				Error: errorMsg.Error(),
			})
		}
	}

	return c.JSON(http.StatusOK, dtos.Response{
		Data: dtos.AuthResponse{
			ExpiresAt: expiresAt,
			JWT:       jwtToken,
			User: dtos.UserInfo{
				ID:    strconv.FormatUint(uint64(user.ID), 10),
				Email: user.Email,
				Name:  user.Name,
			},
		},
		Error: nil,
	})
}
