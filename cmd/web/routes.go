package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	// "github.com/justinas/alice"
)

func (app *Application) routes(mux *http.ServeMux) http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	// dynamic := alice.New(app.SessionManager.LoadAndSave)

	// router.Handler(http.MethodGet, "/", http.HandlerFunc(app.home))
	router.Handler(http.MethodGet, "/", http.HandlerFunc(app.loginSpotify))
	router.Handler(http.MethodGet, "/callback", http.HandlerFunc(app.callbackSpotify))
	// router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	// router.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	// router.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))
	// standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)
	// return standard.Then(router)
	return router
}
