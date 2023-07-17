package middleware

import (
	// "fmt"
	"context"
	"encoding/json"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

type Resource struct {
	Config *oauth2.Config
}

func (r Resource) GenerateYoutubeClient(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var authHeaderType *oauth2.Token
		authHeader := c.Request().Header.Get("AuthorizationYoutube")
		json.Unmarshal([]byte(authHeader), &authHeaderType)
		ctx := context.Background()

		token := &oauth2.Token{
			AccessToken:  authHeaderType.AccessToken,
			TokenType:    "Bearer",
			RefreshToken: authHeaderType.RefreshToken,
			Expiry:       time.Now().Add(time.Hour),
		}

		client := r.Config.Client(ctx, token)

		c.Set("client", client)

		return next(c)
	}
}