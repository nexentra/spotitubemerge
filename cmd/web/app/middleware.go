package app

import (
	// "fmt"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, AuthorizationYoutube, AuthorizationSpotify, X-Requested-With")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

func (app *Application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.InfoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method,
			r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func (app *Application) generateYoutubeClient(next echo.HandlerFunc) echo.HandlerFunc {
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

		client := app.Youtube.Config.Client(ctx, token)

		c.Set("client", client)

		return next(c)
	}
}

// func (app *Application) recoverPanic(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer func() {
// 			if err := recover(); err != nil {
// 				w.Header().Set("Connection", "close")
// 				app.serverError(w, fmt.Errorf("%s", err))
// 			}
// 		}()
// 		next.ServeHTTP(w, r)
// 	})
// }
