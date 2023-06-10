package main

import (
	// "log"
	"net/http"
	// "runtime/pprof"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// "github.com/justinas/alice"
	"github.com/nexentra/spotitubemerge/ui"
)
func (app *Application) routes(mux *http.ServeMux) http.Handler {
	router := echo.New()
	// router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	app.notFound(w)
	// })

	// router.Handler(http.MethodGet, "/", http.HandlerFunc(app.home))
	router.GET("/hello.json", handleHello)
	router.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: frontend.BuildHTTPFS(),
		HTML5:      true,
	}))
	router.GET("/api", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "hello from the echo server",
		})
	})
	// router.GET("/auth/youtube", http.HandlerFunc(app.loginYoutube))
	// router.GET("/auth/youtube/callback", http.HandlerFunc(app.callbackYoutube))
	// router.GET("/youtube-playlist", http.HandlerFunc(app.getYoutubePlaylist))

	// router.GET("/auth/spotify", http.HandlerFunc(app.loginSpotify))
	// router.GET("/auth/spotify/callback", http.HandlerFunc(app.callbackSpotify))
	// router.GET("/spotify-playlist", http.HandlerFunc(app.getSpotifyPlaylist))
	// standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// return standard.Then(router)
	return router
}



func handleHello(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "hello from the echo server",
	})
}