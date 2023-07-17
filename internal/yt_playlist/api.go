package ytplaylist

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
	"github.com/nexentra/spotitubemerge/internal/middleware"
)

func RegisterHandlers(r *echo.Group, config *oauth2.Config, state string, errorLog *log.Logger, infoLog *log.Logger, env map[string]string) {
	res := Resource{
		Config: config,
		State: state,
		ErrorLog: errorLog,
		InfoLog: infoLog,
		Env: env,
	}

	resConfig := middleware.Resource{
		Config: config,
	}

	r.GET("/auth/youtube", res.loginYoutube)
	r.POST("/auth/youtube/callback", res.callbackYoutube)
	r.GET("/youtube-playlist", res.getYoutubePlaylist, resConfig.GenerateYoutubeClient)
	r.GET("/youtube-items", res.getYoutubeItems, resConfig.GenerateYoutubeClient)
}

type Resource struct {
	Config *oauth2.Config
	State string
	Token *oauth2.Token
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Env           map[string]string
}

func (r Resource) loginYoutube(c echo.Context) error {
	authURL := r.Config.AuthCodeURL(r.State, oauth2.AccessTypeOffline)
	fmt.Println("Auth URL: ", authURL)
	// c.Redirect(http.StatusMovedPermanently, authURL)
	return c.JSON(http.StatusOK, echo.Map{
		"authUrl": authURL,
	})
}

type YoutubeCode struct {
	Code string `json:"code"`
}

func (r Resource) callbackYoutube(c echo.Context) error {
	youtubeCode := YoutubeCode{}
	if err := c.Bind(&youtubeCode); err != nil {
		r.ErrorLog.Printf("Failed to bind code: %v", err)
	}
	tok, err := r.Config.Exchange(context.Background(), youtubeCode.Code)
	if err != nil {
		r.ErrorLog.Printf("Unable to retrieve token from web %v", err)
	}

	// Save the token in the variable instead of caching
	r.Token = tok

	fmt.Println("Token: ", r.Token)
	// http.Redirect(w, r, "/playlists", http.StatusFound)
	return c.JSON(http.StatusOK, echo.Map{
		"token": r.Token,
	})
}

func (r Resource) getYoutubePlaylist(c echo.Context) error {
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

func (r Resource) getYoutubeItems(c echo.Context) error {
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







func (r Resource) createPlaylistHandler(c echo.Context)error {
	fmt.Println("Create Playlist")

	client := c.Get("client").(*http.Client)

	playlist, err := CreatePlaylist(client, "My Private Playlist 2")
	if err != nil {
		log.Fatalf("Error creating playlist: %v", err)
	}

	fmt.Println(playlist, "playlist")

	err = AddVideoToPlaylist(client, playlist.Id, "un6ZyFkqFKo")
	if err != nil {
		log.Fatalf("Error adding video to playlist: %v", err)
	}

	err = AddVideoToPlaylist(client, playlist.Id, "mWi9SGKRpys")
	if err != nil {
		log.Fatalf("Error adding video to playlist: %v", err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"playlist": "playlist",
	})
	
}




func AddVideoToPlaylist(client *http.Client, playlistID string, videoID string) error {
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



func CreatePlaylist(client *http.Client, title string) (*youtube.Playlist, error) {
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


func SearchVideoInYoutube(client *http.Client, artist, title string) (string, error) {
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
