package main

import (
	"context"
	"fmt"
	"net/http"

	spotify "github.com/zmb3/spotify/v2"
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
	url := app.Spotify.Authenticator.AuthURL(app.Spotify.State)
	fmt.Println(url)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)
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

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		app.ErrorLog.Print(err)
		return
	}
	fmt.Println("You are logged in as:", user)

	// Get user's playlists
	playlists, err := client.GetPlaylistsForUser(context.Background(), user.ID)
	if err != nil {
		app.ErrorLog.Print(err)
		return
	}
	fmt.Println("Playlists:")
	for _, playlist := range playlists.Playlists {
		fmt.Println("  ", playlist.Name)
	}
}
