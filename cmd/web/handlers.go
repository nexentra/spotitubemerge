package main

import (
	"context"
	"encoding/json"
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
	})
}

func (app *Application) getSpotifyPlaylist(c echo.Context) error {
	var authHeaderType *oauth2.Token
	authHeader := c.Request().Header.Get("Authorization")
	json.Unmarshal([]byte(authHeader), &authHeaderType)
	client := spotify.New(app.Spotify.Authenticator.Client(c.Request().Context(), authHeaderType))

	user, err := client.CurrentUser(c.Request().Context())
	if err != nil {
		app.ErrorLog.Print(err)
	}

	playlists, err := client.GetPlaylistsForUser(c.Request().Context(), user.ID)
	if err != nil {
		app.ErrorLog.Print(err)
	}
	fmt.Println("Playlists:", playlists)
	for _, playlist := range playlists.Playlists {
		fmt.Println("  ", playlist.Name)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"playlists": playlists,
	})
}

func (app *Application) loginYoutube(c echo.Context) error {
	authURL := app.Youtube.Config.AuthCodeURL(app.Youtube.State, oauth2.AccessTypeOnline)
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
	var authHeaderType *oauth2.Token
	authHeader := c.Request().Header.Get("Authorization")
	json.Unmarshal([]byte(authHeader), &authHeaderType)
	url := "https://www.googleapis.com/youtube/v3/playlists?part=snippet&contentDetails&maxResults=" + "25" + "&mine=true&key=" + app.Youtube.Config.ClientID + "&access_token=" + authHeaderType.AccessToken
	fmt.Println("URL: ", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	app.InfoLog.Println(string(responseData))

	return c.JSON(http.StatusOK, echo.Map{
		"playlists": string(responseData),
	})
}

type YoutubePlaylists struct {
	Items []string `json:"items"`
}

func (app *Application) getYoutubeItems(c echo.Context) error {
	var authHeaderType *oauth2.Token
	authHeader := c.Request().Header.Get("Authorization")
	json.Unmarshal([]byte(authHeader), &authHeaderType)

	strings := c.QueryParam("strings")
	fmt.Println("strings: ", strings)

	url := "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&contentDetails&playlistId=" + strings + "&maxResults=" + "25" + "&key=" + app.Youtube.Config.ClientID + "&access_token=" + authHeaderType.AccessToken

	fmt.Println("URL: ", url)
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err.Error())
	}

	client := spotify.New(app.Spotify.Authenticator.Client(c.Request().Context(), authHeaderType))
	results, err := client.Search(c.Request().Context(), "holiday", spotify.SearchTypePlaylist|spotify.SearchTypeAlbum)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("results:", results)

	// handle album results
	if results.Albums != nil {
		fmt.Println("Albums:")
		for _, item := range results.Albums.Albums {
			fmt.Println("   ", item.Name)
		}
	}
	// handle playlist results
	if results.Playlists != nil {
		fmt.Println("Playlists:")
		for _, item := range results.Playlists.Playlists {
			fmt.Println("   ", item.Name)
		}
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	fmt.Println("Response: ", string(responseData))

	return c.JSON(http.StatusOK, echo.Map{
		"items": []string{string(responseData), "item2"},
	})
}
