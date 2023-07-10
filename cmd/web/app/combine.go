package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	spotify "github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
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

	//get spotify user
	user, err := client.CurrentUser(c.Request().Context())
	if err != nil {
		app.ErrorLog.Print(err)
	}

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
	for i, title := range allSpotifyTitles {
		fmt.Println("index: ", i)
		fmt.Println("title: ", title.Track.Artists[0].Name)
		fmt.Println("title: ", title.Track.Name)
	}

	//get youtube items
	ytClient := c.Get("client").(*http.Client)
	service, err := youtube.New(ytClient)
	if err != nil {
		return err
	}

	// Get YouTube items
	var allYoutubeTitles []string
	for _, playlist := range mergerList.YoutubePlaylists {
		call := service.PlaylistItems.List([]string{"snippet"}).
			PlaylistId(playlist).
			MaxResults(50)

		var nextPageToken string
		for {
			if nextPageToken != "" {
				call.PageToken(nextPageToken)
			}

			response, err := call.Do()
			if err != nil {
				fmt.Print(err.Error())
				break
			}

			for _, item := range response.Items {
				title := item.Snippet.Title
				allYoutubeTitles = append(allYoutubeTitles, title)
			}

			nextPageToken = response.NextPageToken
			if nextPageToken == "" {
				break
			}
		}
	}

	// Create a new spotify playlist
	newPlaylist, err := client.CreatePlaylistForUser(context.Background(), user.ID, "testPlaylist", "New playlist for searched tracks", false, false)
	if err != nil {
		fmt.Println("Error creating playlist:", err)
	}

	// Add tracks to the spotify playlist
	for i, title := range allYoutubeTitles {
		fmt.Println("index: ", i)
		fmt.Println("title: ", title)

		// Search for the video in Spotify
		results, err := client.Search(context.Background(), title, spotify.SearchTypeTrack)
		if err != nil {
			fmt.Println("Error searching for video:", err)
			continue
		}

		// Get the first track from the search results
		if len(results.Tracks.Tracks) > 0 {
			firstTrack := results.Tracks.Tracks[0]

			// Add the track to the newly created playlist
			_, err := client.AddTracksToPlaylist(context.Background(), newPlaylist.ID, firstTrack.ID)
			if err != nil {
				fmt.Println("Error adding track to playlist:", err)
			} else {
				fmt.Println("Track added to playlist:", firstTrack.Name)
			}
		} else {
			fmt.Println("No matching track found in Spotify")
		}
	}

	for i, title := range allSpotifyTitles {
		fmt.Println("index: ", i)
		fmt.Println("title: ", title.Track.Artists[0].Name)

			// Add the track to the newly created playlist
			_, err := client.AddTracksToPlaylist(context.Background(), newPlaylist.ID, title.Track.ID)
			if err != nil {
				fmt.Println("Error adding track to playlist:", err)
			} else {
				fmt.Println("Track added to playlist:", title.Track.Name)
			}
		
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
