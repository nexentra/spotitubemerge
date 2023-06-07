package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/nexentra/spotitubemerge/internal/models"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type Application struct {
	ErrorLog             *log.Logger
	InfoLog              *log.Logger
	Spotify              *models.SpotifyModel
}

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

func main() {
	authenticator := spotifyauth.New(
		spotifyauth.WithRedirectURL("http://localhost:8080/callback"),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
			spotifyauth.ScopePlaylistReadPrivate,
		),
		spotifyauth.WithClientID("28489fd2f52440fa90a7191fab27a787"),
		spotifyauth.WithClientSecret("aa9f16eae0eb4d2a89a9e7a8e150e9b3"),
	)

	app := &Application{
		ErrorLog: log.New(log.Writer(), "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		InfoLog:  log.New(log.Writer(), "INFO\t", log.Ldate|log.Ltime),
		Spotify: &models.SpotifyModel{
			Authenticator: authenticator,
			State:         "abc123",
		},
	}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:     ":" + strconv.FormatUint(8080, 10),
		ErrorLog: app.ErrorLog,
		Handler:  app.routes(mux),
	}

	app.InfoLog.Printf("Starting server on http://localhost:%d", 8080)
	err := srv.ListenAndServe()
	app.ErrorLog.Fatal(err)
}
