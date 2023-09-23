package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nexentra/spotitubemerge/ui"
)

func EchoMiddleware(r *echo.Echo, prod string) {
	r.Use(middleware.Logger())
	r.Use(middleware.Recover())
	r.Use(secureHeaders)

	// Middleware for bundling frontend into binary
	if prod  == "production" {
		r.Use(middleware.StaticWithConfig(middleware.StaticConfig{
			Filesystem: frontend.BuildHTTPFS(),
			HTML5:      true,
		}))
	}

	// router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	// 	AllowOrigins: []string{"http://localhost:8080", "https://spotitubemerge.fly.dev"},
	// 	AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodOptions},
	// }))
	// router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	app.notFound(w)
	// })
}

func secureHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// c.Response().Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		c.Response().Header().Set("Content-Security-Policy", "default-src *")
		c.Response().Header().Set("Referrer-Policy", "origin-when-cross-origin")
		c.Response().Header().Set("Access-Control-Allow-Origin", "*")
		c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, AuthorizationYoutube, AuthorizationSpotify, X-Requested-With")
		c.Response().Header().Set("X-Content-Type-Options", "nosniff")
		c.Response().Header().Set("X-Frame-Options", "deny")
		c.Response().Header().Set("X-XSS-Protection", "0")
		return next(c)
	}
}
