package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
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

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	// fmt.Println("Response: ", string(responseData))

	return c.JSON(http.StatusOK, echo.Map{
		"items": []string{string(responseData)},
	})
}
