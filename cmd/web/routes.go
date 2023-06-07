package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *Application) routes(mux *http.ServeMux) http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	// router.Handler(http.MethodGet, "/", http.HandlerFunc(app.home))
	router.Handler(http.MethodGet, "/auth/youtube", http.HandlerFunc(app.loginYoutube))
	router.Handler(http.MethodGet, "/auth/youtube/callback", http.HandlerFunc(app.callbackYoutube))
	router.Handler(http.MethodGet, "/youtube-playlist", http.HandlerFunc(app.getYoutubePlaylist))

	router.Handler(http.MethodGet, "/auth/spotify", http.HandlerFunc(app.loginSpotify))
	router.Handler(http.MethodGet, "/auth/spotify/callback", http.HandlerFunc(app.callbackSpotify))
	router.Handler(http.MethodGet, "/spotify-playlist", http.HandlerFunc(app.getSpotifyPlaylist))
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	return standard.Then(router)
}
