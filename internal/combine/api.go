package combine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nexentra/spotitubemerge/internal/middleware"
	"github.com/nexentra/spotitubemerge/internal/yt_playlist"
	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

func RegisterHandlers(r *echo.Group, authenticator *spotifyauth.Authenticator, config *oauth2.Config, errorLog *log.Logger, infoLog *log.Logger, env map[string]string) {
	res := Resource{
		Authenticator: authenticator,
		ErrorLog: errorLog,
		InfoLog: infoLog,
		Env: env,
	}

	resConfig := middleware.Resource{
		Config: config,
	}

	r.POST("/merge-yt-spotify", res.mergeYtSpotify, resConfig.GenerateYoutubeClient)
}


type MergerList struct {
	YoutubePlaylists []string `json:"youtube-playlists"`
	SpotifyPlaylists []string `json:"spotify-playlists"`
}

type Resource struct {
	Authenticator *spotifyauth.Authenticator
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Env           map[string]string
}


func (r Resource) mergeYtSpotify(c echo.Context) error {
	//get data
	mergerList := MergerList{}
	if err := c.Bind(&mergerList); err != nil {
		r.ErrorLog.Printf("Failed to bind code: %v", err)
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
	client := spotify.New(r.Authenticator.Client(c.Request().Context(), authHeaderTypeSpotify))

	//get spotify user
	user, err := client.CurrentUser(c.Request().Context())
	if err != nil {
		r.ErrorLog.Print(err)
	}

	//get spotify items
	var allSpotifyTitles []spotify.PlaylistItemTrack
	for _, playlist := range mergerList.SpotifyPlaylists {
		playlistItems, err := client.GetPlaylistItems(c.Request().Context(), spotify.ID(playlist))
		if err != nil {
			r.ErrorLog.Print(err)
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

	ytClient := c.Get("client").(*http.Client)
	service, err := youtube.New(ytClient)
	if err != nil {
		return err
	}

	// Get YouTube items
	type AllYoutubeItems struct {
		Title string `json:"title"`
		Id    string `json:"id"`
	}
	var allYoutubeTitles []AllYoutubeItems
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
				allYoutubeTitles = append(allYoutubeTitles, AllYoutubeItems{Title: title, Id: item.Snippet.ResourceId.VideoId})
			}

			nextPageToken = response.NextPageToken
			if nextPageToken == "" {
				break
			}
		}
	}

	// Create a new spotify playlist
	newSpotifyPlaylist, err := client.CreatePlaylistForUser(context.Background(), user.ID, "spotitubeMergePlaylist", "New playlist for searched tracks", false, false)
	if err != nil {
		fmt.Println("Error creating playlist:", err)
	}

	newYoutubePlaylist, err := ytplaylist.CreatePlaylist(ytClient, "spotitubeMergePlaylist")
	if err != nil {
		fmt.Println("Error creating playlist: ", err)
	}

	// Add tracks to the spotify playlist
	for i, items := range allYoutubeTitles {
		fmt.Println("index: ", i)
		fmt.Println("title: ", items.Title)

		// Search for the video in Spotify
		results, err := client.Search(context.Background(), items.Title, spotify.SearchTypeTrack)
		if err != nil {
			fmt.Println("Error searching for video:", err)
			continue
		}

		// Get the first track from the search results
		if len(results.Tracks.Tracks) > 0 {
			firstTrack := results.Tracks.Tracks[0]

			// Add the track to the spotify playlist
			_, err := client.AddTracksToPlaylist(context.Background(), newSpotifyPlaylist.ID, firstTrack.ID)
			if err != nil {
				fmt.Println("Error adding track to playlist:", err)
			} else {
				fmt.Println("Track added to playlist:", firstTrack.Name)
			}
		} else {
			fmt.Println("No matching track found in Spotify")
		}

		// Add the video to the youtube playlist
		err = ytplaylist.AddVideoToPlaylist(ytClient, newYoutubePlaylist.Id, items.Id)
		if err != nil {
			fmt.Println("Error adding video to playlist: ", err)
		} else {
			fmt.Println("Track added to playlist:", items.Title)
		}
	}

	for i, title := range allSpotifyTitles {
		fmt.Println("index: ", i)
		fmt.Println("title: ", title.Track.Artists[0].Name)
		// Add the track to the newly created playlist
		_, err := client.AddTracksToPlaylist(context.Background(), newSpotifyPlaylist.ID, title.Track.ID)
		if err != nil {
			fmt.Println("Error adding track to playlist:", err)
		} else {
			fmt.Println("Track added to playlist:", title.Track.Name)
		}

		// Search for the track in YouTube
		videoID, err := ytplaylist.SearchVideoInYoutube(ytClient, title.Track.Artists[0].Name, title.Track.Name)
		if err != nil {
			fmt.Println("Error searching for track in YouTube:", err)
			continue
		}

		// Add the video to the YouTube playlist
		err = ytplaylist.AddVideoToPlaylist(ytClient, newYoutubePlaylist.Id, videoID)
		if err != nil {
			fmt.Println("Error adding video to YouTube playlist: ", err)
		} else {
			fmt.Println("Video added to YouTube playlist:", title.Track.Name)
		}

	}

	return c.JSON(http.StatusOK, echo.Map{
		"items": allYoutubeTitles,
	})
}
