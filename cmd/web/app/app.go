package app

import (
	"log"

	"github.com/nexentra/spotitubemerge/internal/models"
)

type Application struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Spotify  *models.SpotifyModel
	Youtube  *models.YoutubeModel
	Env      map[string]string
}