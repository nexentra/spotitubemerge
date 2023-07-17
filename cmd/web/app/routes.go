package app

import (
	"net/http"
	// "runtime/pprof"

	// "github.com/justinas/alice"
	// "github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nexentra/spotitubemerge/internal/combine"
	spotifyplaylist "github.com/nexentra/spotitubemerge/internal/spotify_playlist"
	ytplaylist "github.com/nexentra/spotitubemerge/internal/yt_playlist"
	// "github.com/justinas/alice"
	// "github.com/nexentra/spotitubemerge/ui"
)
func (app *Application) Routes(mux *http.ServeMux) http.Handler {
	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())
	router.Use(secureHeaders)

	router.GET("/devices", getDevices)
	router.POST("/devices", createDevice)
	router.PUT("/devices/:id", upgradeDevice)
	router.GET("/login", login, loginMiddleware)

	// Middleware for bundling frontend into binary
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
	

	spotifyplaylist.RegisterHandlers(apiRoute, app.Spotify.Authenticator, app.Spotify.RedirectURI, app.Spotify.State, app.ErrorLog, app.InfoLog, app.Env)
	ytplaylist.RegisterHandlers(apiRoute, app.Youtube.Config, app.Youtube.State, app.ErrorLog, app.InfoLog, app.Env)
	combine.RegisterHandlers(apiRoute,app.Spotify.Authenticator,  app.Youtube.Config, app.ErrorLog, app.InfoLog, app.Env)

	// apiRoute.POST("/mytest", app.createPlaylistHandler, app.generateYoutubeClient)
	// standard := alice.New(app.logRequest, secureHeaders)
	// return standard.Then(router)
	return router
}


func secureHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		c.Response().Header().Set("Referrer-Policy", "origin-when-cross-origin")
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, AuthorizationYoutube, AuthorizationSpotify, X-Requested-With")
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		c.Response().Header().Set("X-Frame-Options", "deny")
		c.Response().Header().Set("X-XSS-Protection", "0")
		return next(c)
	}
}