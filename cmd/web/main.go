package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/nexentra/spotitubemerge/internal/models"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Spotify  *models.SpotifyModel
	Youtube  *models.YoutubeModel
	Env      map[string]string
}

type configType struct {
	logToFile bool
}

func main() {
	envFile, _ := godotenv.Read(".env")
	var cfg configType
	var errorLog *log.Logger
	var infoLog *log.Logger

	flag.BoolVar(&cfg.logToFile, "log", false, "Enable logging to file")
	flag.Parse()

	if cfg.logToFile {
		infoFile, err := os.OpenFile("tmp/info.log", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}

		err = infoFile.Truncate(0)
		if err != nil {
			log.Fatal(err)
		}

		errFile, err := os.OpenFile("tmp/error.log", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			log.Fatal(err)
		}

		err = errFile.Truncate(0)
		if err != nil {
			log.Fatal(err)
		}

		defer infoFile.Close()
		defer errFile.Close()
		infoLog = log.New(infoFile, "INFO\t", log.Ldate|log.Ltime)
		errorLog = log.New(errFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	} else {
		infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
		errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	}

	var ytConfigFile string

	if envFile["NODE_ENV"] == "production" {
		ytConfigFile = "client_secret.json"
	} else {
		ytConfigFile = "dev_client.json"
	}

	b, err := ioutil.ReadFile(ytConfigFile)
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, youtube.YoutubeScope)
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v", err)
	}

	var redirectUri string

	if envFile["NODE_ENV"] == "production" {
		redirectUri = "https://spotitubemerge.nexentra.online/auth/spotify/callback"
	} else {
		redirectUri = "http://localhost:8080/auth/spotify/callback"
	}

	authenticator := spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectUri),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopePlaylistModifyPublic,
			spotifyauth.ScopePlaylistModifyPrivate,
			
		),
		spotifyauth.WithClientID(envFile["SPOTIFY_CLIENT_ID"]),
		spotifyauth.WithClientSecret(envFile["SPOTIFY_CLIENT_SECRET"]),
	)

	application := &Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
		Spotify: &models.SpotifyModel{
			Authenticator: authenticator,
			State:         envFile["SPOTIFY_STATE"],
		},
		Youtube: &models.YoutubeModel{
			Config: config,
			State:  envFile["YOUTUBE_STATE"],
		},
		Env: envFile,
	}

	mux := http.NewServeMux()
	prometheusMux := http.NewServeMux()

	srv := &http.Server{
		Addr:     ":" + strconv.FormatUint(8080, 10),
		ErrorLog: application.ErrorLog,
		Handler:  application.Routes(mux),
	}

	prometheusSrv := &http.Server{
		Addr:     ":" + strconv.FormatUint(8081, 10),
		ErrorLog: application.ErrorLog,
		Handler:  application.PrometheusRoutes(prometheusMux),
	}

go func(){
	application.InfoLog.Printf("Starting server on http://localhost:%d", 8080)
	err = srv.ListenAndServe()
	application.ErrorLog.Fatal(err)
}()

go func(){
	application.InfoLog.Printf("Starting server on http://localhost:%d", 8081)
	err = prometheusSrv.ListenAndServe()
	application.ErrorLog.Fatal(err)
}()

select {}
}
