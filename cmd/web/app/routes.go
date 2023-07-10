package app

import (
	// "log"
	"net/http"
	// "runtime/pprof"

	"github.com/justinas/alice"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	// "github.com/justinas/alice"
	// "github.com/nexentra/spotitubemerge/ui"
)
func (app *Application) Routes(mux *http.ServeMux) http.Handler {
	router := echo.New()
	// router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	// router.Use(middleware.StaticWithConfig(middleware.StaticConfig{
	// 	Filesystem: frontend.BuildHTTPFS(),
	// 	HTML5:      true,
	// }))
	// router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"http://localhost:8080", "https://spotitubemerge.fly.dev"},
	// 	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
	// }))
	// router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	app.notFound(w)
	// })

	apiRoute := router.Group("/api")
	apiRoute.GET("/auth/youtube", app.loginYoutube)
	apiRoute.POST("/auth/youtube/callback", app.callbackYoutube)
	apiRoute.GET("/youtube-playlist", app.getYoutubePlaylist, app.generateYoutubeClient)
	apiRoute.GET("/youtube-items", app.getYoutubeItems, app.generateYoutubeClient)

	apiRoute.GET("/auth/spotify", app.loginSpotify)
	apiRoute.POST("/auth/spotify/callback", app.callbackSpotify)
	apiRoute.GET("/spotify-playlist", app.getSpotifyPlaylist)
	apiRoute.GET("/spotify-items", app.getSpotifyItems)
	apiRoute.GET("/search-spotify-items", app.searchSpotifyItems)

	apiRoute.POST("/merge-yt-spotify", app.mergeYtSpotify,app.generateYoutubeClient)
	apiRoute.POST("/mytest", app.createPlaylistHandler, app.generateYoutubeClient)
	standard := alice.New(app.logRequest, secureHeaders)
	return standard.Then(router)
	// return router
}