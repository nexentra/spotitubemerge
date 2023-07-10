package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "time"

	"github.com/labstack/echo/v4"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

func (app *Application) loginYoutube(c echo.Context) error {
	authURL := app.Youtube.Config.AuthCodeURL(app.Youtube.State, oauth2.AccessTypeOffline)
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
	client := c.Get("client").(*http.Client)
	service, err := youtube.New(client)
	if err != nil {
		return err
	}

	call := service.Playlists.List([]string{"snippet,contentDetails"}).
		MaxResults(50).
		Mine(true)

	var allResponseData []*youtube.Playlist

	for {
		response, err := call.Do()
		if err != nil {
			return err
		}

		allResponseData = append(allResponseData, response.Items...)

		if response.NextPageToken == "" {
			break
		}

		call.PageToken(response.NextPageToken)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"playlists": allResponseData,
	})
}

type YoutubePlaylists struct {
	Items []string `json:"items"`
}

func (app *Application) getYoutubeItems(c echo.Context) error {
	var authHeaderType *oauth2.Token
	authHeader := c.Request().Header.Get("Authorization")
	json.Unmarshal([]byte(authHeader), &authHeaderType)

	playlistID := c.QueryParam("strings")
	fmt.Println("playlistID: ", playlistID)

	client := c.Get("client").(*http.Client)
	service, err := youtube.New(client)
	if err != nil {
		return err
	}

	call := service.PlaylistItems.List([]string{"snippet,contentDetails"}).
		PlaylistId(playlistID).
		MaxResults(25)

	response, err := call.Do()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"items": response.Items,
	})
}







func (app *Application) createPlaylistHandler(c echo.Context)error {
	fmt.Println("Create Playlist")

	client := c.Get("client").(*http.Client)

	playlist, err := createPlaylist(client, "My Private Playlist 2")
	if err != nil {
		log.Fatalf("Error creating playlist: %v", err)
	}

	fmt.Println(playlist, "playlist")

	err = addVideoToPlaylist(client, playlist.Id, "un6ZyFkqFKo")
	if err != nil {
		log.Fatalf("Error adding video to playlist: %v", err)
	}

	err = addVideoToPlaylist(client, playlist.Id, "mWi9SGKRpys")
	if err != nil {
		log.Fatalf("Error adding video to playlist: %v", err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"playlist": "playlist",
	})
	
}




func addVideoToPlaylist(client *http.Client, playlistID string, videoID string) error {
	service, err := youtube.New(client)
	if err != nil {
		return err
	}

	playlistItem := &youtube.PlaylistItem{
		Snippet: &youtube.PlaylistItemSnippet{
			PlaylistId: playlistID,
			ResourceId: &youtube.ResourceId{
				Kind:    "youtube#video",
				VideoId: videoID,
			},
		},
	}

	call := service.PlaylistItems.Insert([]string{"snippet"}, playlistItem)
	_, err = call.Do()
	if err != nil {
		return err
	}

	return nil
}



func createPlaylist(client *http.Client, title string) (*youtube.Playlist, error) {
	service, err := youtube.New(client)
	if err != nil {
		return nil, err
	}

	playlist := &youtube.Playlist{
		Snippet: &youtube.PlaylistSnippet{
			Title:       title,
			Description: "Private playlist created using the YouTube API v3",
		},
		Status: &youtube.PlaylistStatus{
			PrivacyStatus: "private",
		},
	}

	call := service.Playlists.Insert([]string{"snippet,status"}, playlist)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	return response, nil
}


func searchVideoInYoutube(client *http.Client, artist, title string) (string, error) {
	service, err := youtube.New(client)
	if err != nil {
		return "", err
	}

	query := artist + " " + title
	call := service.Search.List([]string{"id"}).Q(query).Type("video").MaxResults(1)
	response, err := call.Do()
	if err != nil {
		return "", err
	}

	if len(response.Items) == 0 {
		return "", fmt.Errorf("No matching video found in YouTube")
	}

	videoID := response.Items[0].Id.VideoId
	return videoID, nil
}
