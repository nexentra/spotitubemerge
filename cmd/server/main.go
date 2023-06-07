package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/zmb3/spotify/v2"
	"github.com/zmb3/spotify/v2/auth"
)

type Application struct {
	spotifyAuthenticator *spotifyauth.Authenticator
	spotifyClient        *spotify.Client
	spotifyRedirectURI   string
	spotifyState         string
	spotifyID            string
	spotifySecret        string
}

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

const redirectURI = "http://localhost:8080/callback"

var (
	spotifyID     = "28489fd2f52440fa90a7191fab27a787"
	spotifySecret = "aa9f16eae0eb4d2a89a9e7a8e150e9b3"
	authenticator = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserReadPrivate,
			spotifyauth.ScopePlaylistReadCollaborative,
			spotifyauth.ScopePlaylistReadPrivate,
		),
		spotifyauth.WithClientID(spotifyID),
		spotifyauth.WithClientSecret(spotifySecret),
	)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

func main() {
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := authenticator.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	client := <-ch

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := authenticator.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	client := spotify.New(authenticator.Client(r.Context(), tok))
	fmt.Fprint(w, "Login Completed!")

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user)

	// Get user's playlists
	playlists, err := client.GetPlaylistsForUser(context.Background(), user.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Playlists:")
	for _, playlist := range playlists.Playlists {
		fmt.Println("  ", playlist.Name)
	}
	ch <- client
}
