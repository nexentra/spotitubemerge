package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	spotify "github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

// func (app *Application) home(c echo.Context) error {
// 	// snippets, err := app.Snippets.Latest()
// 	// if err != nil {
// 	// 	app.serverError(w, err)
// 	// 	return
// 	// }
// 	// data := app.newTemplateData(r)
// 	// data.Snippets = snippets

// 	// app.render(w, http.StatusOK, "home.html", data)

// }

func (app *Application) loginSpotify(c echo.Context) error {
	authURL := app.Spotify.Authenticator.AuthURL(app.Spotify.State)
	fmt.Println("Auth URL: ", authURL)
	return c.JSON(http.StatusOK, echo.Map{
		"authUrl": authURL,
	})
}


type SpotifyCode struct {
	Code string `json:"code"`
}

func (app *Application) callbackSpotify(c echo.Context) error {
	spotifyCode := SpotifyCode{}
	if err := c.Bind(&spotifyCode); err != nil {
		app.ErrorLog.Printf("Failed to bind code: %v", err)
	}
	fmt.Println("Code: ", spotifyCode.Code)

	// state := app.Spotify.State
	tok, err := app.Spotify.Authenticator.Exchange(c.Request().Context(), spotifyCode.Code)
	if err != nil {
		app.ErrorLog.Print(err)
	}
	fmt.Println("Token: ", tok)

	// if st := c.FormValue("state"); st != state {
	// 	app.ErrorLog.Printf("State mismatch: %s != %s\n", st, state)
	// }

	client := spotify.New(app.Spotify.Authenticator.Client(c.Request().Context(), tok))
	fmt.Println("Client: ", client)
	app.Spotify.Client = client

	user, err := app.Spotify.Client.CurrentUser(c.Request().Context())
	if err != nil {
		app.ErrorLog.Print(err)
	}
	app.Spotify.UserId = user.ID
	fmt.Println("You are logged in as:", user.ID)
	return c.JSON(http.StatusOK, echo.Map{
		"token": tok,
		"client": client,
	})
}

func (app *Application) getSpotifyPlaylist(c echo.Context) error {
	// Get user's playlists
	playlists, err := app.Spotify.Client.GetPlaylistsForUser(context.Background(), app.Spotify.UserId)
	if err != nil {
		app.ErrorLog.Print(err)
	}
	fmt.Println("Playlists:")
	for _, playlist := range playlists.Playlists {
		fmt.Println("  ", playlist.Name)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"playlists": playlists,
	})
}

func (app *Application) loginYoutube(c echo.Context) error {
	authURL := app.Youtube.Config.AuthCodeURL(app.Spotify.State, oauth2.AccessTypeOnline)
	fmt.Println("Auth URL: ", authURL)
	// c.Redirect(http.StatusMovedPermanently, authURL)
	return c.JSON(http.StatusOK, echo.Map{
		"authUrl": authURL,
	})
}

type YoutubeCode struct {
	Code string `json:"code"`
}

func (app *Application) callbackYoutube(c echo.Context) error {
	youtubeCode := YoutubeCode{}
	if err := c.Bind(&youtubeCode); err != nil {
		app.ErrorLog.Printf("Failed to bind code: %v", err)
	}
	tok, err := app.Youtube.Config.Exchange(context.Background(), youtubeCode.Code)
	if err != nil {
		app.ErrorLog.Printf("Unable to retrieve token from web %v", err)
	}

	// Save the token in the variable instead of caching
	app.Youtube.Token = tok

	fmt.Println("Token: ", app.Youtube.Token)
	// http.Redirect(w, r, "/playlists", http.StatusFound)
	return c.JSON(http.StatusOK, echo.Map{
		"token": app.Youtube.Token,
	})
}

func (app *Application) getYoutubePlaylist(c echo.Context) error {
	fmt.Println("Token: ", app.Youtube.Token.AccessToken)
	url := "https://www.googleapis.com/youtube/v3/playlists?part=snippet%2CcontentDetails&maxResults=" + "25" + "&mine=true&key=" + app.Youtube.Config.ClientID + "&access_token=" + app.Youtube.Token.AccessToken

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

	return c.JSON(http.StatusOK, echo.Map{
		"playlists": body,
	})
}
