package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	spotify "github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

type MergerList struct {
	YoutubePlaylists []string `json:"youtube-playlists"`
	SpotifyPlaylists []string `json:"spotify-playlists"`
}

func (app *Application) mergeYtSpotify(c echo.Context) error {
	//get data
	mergerList := MergerList{}
	if err := c.Bind(&mergerList); err != nil {
		app.ErrorLog.Printf("Failed to bind code: %v", err)
	}

	//get auth headers items
	var authHeaderTypeYoutube *oauth2.Token
	var authHeaderTypeSpotify *oauth2.Token
	authHeaderYoutube := c.Request().Header.Get("AuthorizationYoutube")
	authHeaderSpotify := c.Request().Header.Get("AuthorizationSpotify")
	json.Unmarshal([]byte(authHeaderSpotify), &authHeaderTypeSpotify)
	json.Unmarshal([]byte(authHeaderYoutube), &authHeaderTypeYoutube)
	if authHeaderTypeSpotify == nil || authHeaderTypeYoutube == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "no auth header",
		})
	}

	//spotify client
	client := spotify.New(app.Spotify.Authenticator.Client(c.Request().Context(), authHeaderTypeSpotify))

	//get spotify items
	var allSpotifyTitles []spotify.PlaylistItemTrack
	for _, playlist := range mergerList.SpotifyPlaylists {
		playlistItems, err := client.GetPlaylistItems(c.Request().Context(), spotify.ID(playlist))
		if err != nil {
			app.ErrorLog.Print(err)
		}
		for _, playlistItems := range playlistItems.Items {
			allSpotifyTitles = append(allSpotifyTitles, playlistItems.Track)
			fmt.Println("  ", playlistItems.Track)
		}
	}

	fmt.Println("allSpotifyTitles: ", allSpotifyTitles)
	for i, title := range allSpotifyTitles{
		fmt.Println("index: ", i)
		fmt.Println("title: ", title.Track.Artists[0].Name)
		fmt.Println("title: ", title.Track.Name)
	}

	//get youtube items
	var allYoutubeTitles []string
	var nextPageToken string
	for _, playlist := range mergerList.YoutubePlaylists {
		for {
			url := "https://www.googleapis.com/youtube/v3/playlistItems?part=snippet&playlistId=" + playlist + "&maxResults=50" + "&key=" + app.Youtube.Config.ClientID + "&access_token=" + authHeaderTypeYoutube.AccessToken
			if nextPageToken != "" {
				url += "&pageToken=" + nextPageToken
			}

			response, err := http.Get(url)
			if err != nil {
				fmt.Print(err.Error())
				break
			}

			responseData, err := ioutil.ReadAll(response.Body)
			if err != nil {
				fmt.Print(err.Error())
				break
			}

			var data map[string]interface{}
			err = json.Unmarshal(responseData, &data)
			if err != nil {
				fmt.Println("Error parsing JSON response:", err)
				break
			}

			// Extract the 'items' field as a slice of interface{} from the 'data' map.
			items, ok := data["items"].([]interface{})
			if !ok {
				fmt.Println("Invalid 'items' field in the response")
				break
			}

			// Iterate over each item in the 'items' slice.
			for _, item := range items {
				// Type assert the item as a map[string]interface{}.
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					fmt.Println("Invalid item format in the response")
					continue
				}

				// Access the 'snippet' field as a map[string]interface{}.
				snippet, ok := itemMap["snippet"].(map[string]interface{})
				if !ok {
					fmt.Println("Invalid 'snippet' field in the item")
					continue
				}

				// Access the 'title' field under 'snippet'.
				title, ok := snippet["title"].(string)
				if !ok {
					fmt.Println("Invalid 'title' field in the snippet")
					continue
				}

			allYoutubeTitles = append(allYoutubeTitles, string(title))
			}
			if _, ok := data["nextPageToken"]; ok {
				nextPageToken = data["nextPageToken"].(string)
			} else {
				break
			}
		}
	}

	for i, title := range allYoutubeTitles{
		fmt.Println("index: ", i)
		fmt.Println("title: ", title)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"items": allYoutubeTitles,
	})
}

// func (app *Application) searchSpotifyItems(c echo.Context) error {
// 	var authHeaderType *oauth2.Token
// 	authHeader := c.Request().Header.Get("Authorization")
// 	json.Unmarshal([]byte(authHeader), &authHeaderType)

// 	strings := c.QueryParam("strings")
// 	fmt.Println("strings: ", strings)

// 	// searching tracks with given name
// 	client := spotify.New(app.Spotify.Authenticator.Client(c.Request().Context(), authHeaderType))
// 	results, err := client.Search(c.Request().Context(), strings, spotify.SearchTypeTrack) //spotify.SearchTypePlaylist|spotify.SearchTypeAlbum
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println("results:", results)

// 	// handle songs results
// 	if results.Tracks != nil {
// 		return c.JSON(http.StatusOK, echo.Map{
// 			"items": results.Tracks.Tracks,
// 		})
// 	}

// 	return c.JSON(http.StatusOK, echo.Map{
// 		"items": []string{},
// 	})
// }
