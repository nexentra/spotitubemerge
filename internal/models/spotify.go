package models

import(
	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type SpotifyModel struct {
	Authenticator *spotifyauth.Authenticator
	Client        *spotify.Client
	RedirectURI   string
	State         string
	UserId 	  string
}