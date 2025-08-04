package middlewares

import (
	"backend/repositories"
	"backend/utils"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func TokenRefreshMiddleware(
	cfg *oauth2.Config,
	repo repositories.UserRepository,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var encryptKey64 = os.Getenv("ENCRYPTION_SECRET_KEY")
			key, err := base64.StdEncoding.DecodeString(encryptKey64)

			googleId := (c.Get("google_id")).(string)
			u, err := repo.GetByGoogleID(googleId)
			if err != nil {
				fmt.Println("user not found")
				return echo.NewHTTPError(http.StatusUnauthorized, "user not found")
			}
			accesstokenDecrypted, err := utils.DecryptAccessToken(u.AccessToken, key)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "failed to decrypt access token")
			}

			accessToken := accesstokenDecrypted
			if time.Until(u.Expiry) <= 5*time.Minute {
				tok, err := utils.RefreshAccessToken(c.Request().Context(), cfg, u.RefreshToken)
				if err != nil {
					log.Println("Refresh failed:", err)
					return echo.NewHTTPError(http.StatusUnauthorized, "Token refresh failed")
				} else {
					accessToken = tok.AccessToken

					c.Set("googleAccessToken", accessToken)
					u.Expiry = tok.Expiry
					if err := repo.Update(u); err != nil {
						log.Println("Failed to update user expiry:", err)
					}
				}
			}

			c.Set("googleAccessToken", accessToken)

			return next(c)
		}
	}
}
