package models

import (
	"golang.org/x/oauth2"
)

type YoutubeModel struct {
	Config *oauth2.Config
	State string
	Token *oauth2.Token
	
}