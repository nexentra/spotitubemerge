package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	spotify "github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	// snippets, err := app.Snippets.Latest()
	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	// data := app.newTemplateData(r)
	// data.Snippets = snippets

	// app.render(w, http.StatusOK, "home.html", data)

	fmt.Fprint(w, "Hello, world!")
}

func (app *Application) loginSpotify(w http.ResponseWriter, r *http.Request) {
	authURL := app.Spotify.Authenticator.AuthURL(app.Spotify.State)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (app *Application) callbackSpotify(w http.ResponseWriter, r *http.Request) {
	state := app.Spotify.State
	tok, err := app.Spotify.Authenticator.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		app.ErrorLog.Print(err)
		return
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		app.ErrorLog.Printf("State mismatch: %s != %s\n", st, state)
		return
	}

	client := spotify.New(app.Spotify.Authenticator.Client(r.Context(), tok))
	fmt.Fprint(w, "Login Completed!")

	app.Spotify.Client = client

	user, err := app.Spotify.Client.CurrentUser(context.Background())
	if err != nil {
		app.ErrorLog.Print(err)
		return
	}
	app.Spotify.UserId = user.ID
	fmt.Println("You are logged in as:", user)
}

func (app *Application) getSpotifyPlaylist(w http.ResponseWriter, r *http.Request) {
	// Get user's playlists
	playlists, err := app.Spotify.Client.GetPlaylistsForUser(context.Background(), app.Spotify.UserId)
	if err != nil {
		app.ErrorLog.Print(err)
		return
	}
	fmt.Println("Playlists:")
	for _, playlist := range playlists.Playlists {
		fmt.Println("  ", playlist.Name)
	}
	}


func (app *Application) loginYoutube(w http.ResponseWriter, r *http.Request) {
	authURL := app.Youtube.Config.AuthCodeURL(app.Spotify.State, oauth2.AccessTypeOnline)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func (app *Application) callbackYoutube(w http.ResponseWriter, r *http.Request) {
	var code = r.URL.Query().Get("code")
	tok, err := app.Youtube.Config.Exchange(context.Background(), code)
	if err != nil {
		app.ErrorLog.Printf("Unable to retrieve token from web %v", err)
	}

	// Save the token in the variable instead of caching
	app.Youtube.Token = tok

	fmt.Println("Token: ", app.Youtube.Token)
	// http.Redirect(w, r, "/playlists", http.StatusFound)
}

func (app *Application) getYoutubePlaylist(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Token: ", app.Youtube.Token.AccessToken)
	url := "https://www.googleapis.com/youtube/v3/playlists?part=snippet%2CcontentDetails&maxResults="+"25"+"&mine=true&key="+ app.Youtube.Config.ClientID+ "&access_token=" + app.Youtube.Token.AccessToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		app.ErrorLog.Printf("Failed to create request: %v", err)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		app.ErrorLog.Printf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.ErrorLog.Printf("Failed to read response body: %v", err)
	}

	fmt.Println(string(body))
}