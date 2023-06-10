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
	// router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Filesystem: frontend.BuildHTTPFS(),
		HTML5:      true,
	}))
	// router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"http://localhost:8080", "https://spotitubemerge.fly.dev"},
	// 	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
	// }))
	// router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	app.notFound(w)
	// })

	// router.Handler(http.MethodGet, "/", http.HandlerFunc(app.home))
	router.GET("/hello.json", handleHello)
	router.GET("/api", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"message": "hello from the echo server",
		})
	})
	router.GET("/api/auth/youtube", app.loginYoutube)
	router.POST("/api/auth/youtube/callback", app.callbackYoutube)
	router.GET("/api/youtube-playlist", app.getYoutubePlaylist)

	router.GET("/api/auth/spotify", app.loginSpotify)
	router.POST("/api/auth/spotify/callback", app.callbackSpotify)
	router.GET("/api/spotify-playlist", app.getSpotifyPlaylist)
	// standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// return standard.Then(router)
	return router
}



func handleHello(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{
		"message": "hello from the echo server",
	})
}