package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/nexentra/spotitubemerge/internal/models"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

type Application struct {
	ErrorLog             *log.Logger
	InfoLog              *log.Logger
	Spotify              *models.SpotifyModel
	Youtube              *models.YoutubeModel
	Env 				map[string]string
}

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

func main() {
	envFile, _ := godotenv.Read(".env")
	
	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v", err)
	}
	
	authenticator := spotifyauth.New(
		spotifyauth.WithRedirectURL("http://localhost:8080/auth/spotify/callback"),
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
		Youtube: &models.YoutubeModel{
			Config: config,
			State:  "abc123",
		},
		Env: envFile,
	}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:     ":" + strconv.FormatUint(8080, 10),
		ErrorLog: app.ErrorLog,
		Handler:  app.routes(mux),
	}

	app.InfoLog.Printf("Starting server on http://localhost:%d", 8080)
	err = srv.ListenAndServe()
	app.ErrorLog.Fatal(err)
}
