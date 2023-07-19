package main

import (
	"github.com/labstack/echo/v4"
	"github.com/nexentra/spotitubemerge/internal/combine"
	"github.com/nexentra/spotitubemerge/internal/middleware"
	spotifyplaylist "github.com/nexentra/spotitubemerge/internal/spotify_playlist"
	ytplaylist "github.com/nexentra/spotitubemerge/internal/yt_playlist"
	"net/http"
)

func (app *Application) Routes(mux *http.ServeMux) http.Handler {
	router := echo.New()

	middleware.EchoMiddleware(router, app.Env["PRODUCTION_MODE"])

	router.GET("/devices", getDevices)
	router.POST("/devices", createDevice)
	router.PUT("/devices/:id", upgradeDevice)
	router.GET("/login", login, loginMiddleware)

	apiRoute := router.Group("/api")

	spotifyplaylist.RegisterHandlers(apiRoute, app.Spotify.Authenticator, app.Spotify.RedirectURI, app.Spotify.State, app.ErrorLog, app.InfoLog, app.Env)
	ytplaylist.RegisterHandlers(apiRoute, app.Youtube.Config, app.Youtube.State, app.ErrorLog, app.InfoLog, app.Env)
	combine.RegisterHandlers(apiRoute, app.Spotify.Authenticator, app.Youtube.Config, app.ErrorLog, app.InfoLog, app.Env)

	return router
}
