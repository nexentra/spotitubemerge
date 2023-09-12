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
	"github.com/redis/go-redis/v9"
	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"google.golang.org/api/youtube/v3"
)

func RegisterHandlers(r *echo.Group, authenticator *spotifyauth.Authenticator, config *oauth2.Config, errorLog *log.Logger, infoLog *log.Logger, env map[string]string, redisClient *redis.Client) {
	res := &Resource{
		Authenticator: authenticator,
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		Env:           env,
		RedisClient:   redisClient,
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
	RedisClient   *redis.Client
}

type AllYoutubeItems struct {
	Title string `json:"title"`
	Id    string `json:"id"`
}

func (r *Resource) mergeYtSpotify(c echo.Context) error {
	//get data
	mergerList := MergerList{}
	if err := c.Bind(&mergerList); err != nil {
		r.ErrorLog.Printf("Failed to bind code: %v", err)
	}

	//get auth headers items
	authHeaderTypeSpotify, _, err := getAuthHeaders(c)
	if err != nil {
		r.ErrorLog.Printf("Failed task: %v", err)
	}

	//spotify client
	spotifyClient, user, err := r.createSpotifyClientAndUser(c, authHeaderTypeSpotify)
	if err != nil {
		r.ErrorLog.Printf("Failed task: %v", err)
	}

	//get spotify items
	allSpotifyTitles,err := r.getAllSpotifyItems(c, mergerList, spotifyClient)
	if err != nil {
		r.ErrorLog.Printf("Failed task: %v", err)
	}

	//create youtube client and service
	ytClient, service, err := r.createYoutubeClientAndService(c)
	if err != nil {
		r.ErrorLog.Printf("Failed task: %v", err)
	}


	// Get YouTube items
	allYoutubeTitles,err := r.getAllYoutubeItems(mergerList, service)
	if err != nil {
		r.ErrorLog.Printf("Failed task: %v", err)
	}

	// Create a new spotify playlist
	newSpotifyPlaylist, err := spotifyClient.CreatePlaylistForUser(context.Background(), user.ID, "spotitubeMergePlaylist", "New playlist for searched tracks", false, false)
	if err != nil {
		r.ErrorLog.Printf("Error creating playlist: %v", err)
	}

	// Create a new youtube playlist
	newYoutubePlaylist, err := ytplaylist.CreatePlaylist(ytClient, "spotitubeMergePlaylist")
	if err != nil {
		r.ErrorLog.Printf("Error creating playlist: %v", err)
	}

	// Add tracks to the spotify playlist
	err = r.addTracksToSpotifyPlaylist(allYoutubeTitles,spotifyClient, newSpotifyPlaylist, ytClient,newYoutubePlaylist)
	if err != nil {
		r.ErrorLog.Printf("Failed task: %v", err)
	}

	//Add tracks to the youtube playlist
	err = r.addTracksToYoutubePlaylist(allSpotifyTitles,spotifyClient, newSpotifyPlaylist, ytClient,newYoutubePlaylist)
	if err != nil {
		r.ErrorLog.Printf("Failed task: %v", err)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"items": allYoutubeTitles,
	})
}

func getAuthHeaders(c echo.Context) (*oauth2.Token, *oauth2.Token, error) {
	var authHeaderTypeYoutube *oauth2.Token
	var authHeaderTypeSpotify *oauth2.Token
	authHeaderYoutube := c.Request().Header.Get("AuthorizationYoutube")
	authHeaderSpotify := c.Request().Header.Get("AuthorizationSpotify")
	json.Unmarshal([]byte(authHeaderSpotify), &authHeaderTypeSpotify)
	json.Unmarshal([]byte(authHeaderYoutube), &authHeaderTypeYoutube)
	if authHeaderTypeSpotify == nil || authHeaderTypeYoutube == nil {
		return nil, nil, fmt.Errorf("no auth header")
	}
	return authHeaderTypeSpotify, authHeaderTypeYoutube , nil
}

func (r *Resource) createSpotifyClientAndUser(c echo.Context, authHeaderTypeSpotify *oauth2.Token) (*spotify.Client, *spotify.PrivateUser, error){
	spotifyClient := spotify.New(r.Authenticator.Client(c.Request().Context(), authHeaderTypeSpotify))

	//get spotify user
	user, err := spotifyClient.CurrentUser(c.Request().Context())
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't get user: %w", err)
	}

	return spotifyClient, user, err
}

func (r *Resource) getAllSpotifyItems(c echo.Context, mergerList MergerList, spotifyClient *spotify.Client)([]spotify.PlaylistItemTrack, error){
	var allSpotifyTitles []spotify.PlaylistItemTrack
	for _, playlist := range mergerList.SpotifyPlaylists {
		playlistItems, err := spotifyClient.GetPlaylistItems(c.Request().Context(), spotify.ID(playlist))
		if err != nil {
			return nil, fmt.Errorf("couldn't get playlist items: %w", err)
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

	
	return allSpotifyTitles, nil
}

func (r *Resource) createYoutubeClientAndService(c echo.Context)(*http.Client, *youtube.Service, error){
	ytClient := c.Get("client").(*http.Client)
	service, err := youtube.New(ytClient)
	if err != nil {
		return nil, nil, fmt.Errorf("Error creating new YouTube client: %v", err)
	}

	return ytClient, service, err
}

func (r *Resource) getAllYoutubeItems(mergerList MergerList, service *youtube.Service)([]AllYoutubeItems, error){
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
				return nil, fmt.Errorf("Error fetching playlist items: %v", err)
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

	return allYoutubeTitles, nil
}

func (r *Resource) addTracksToSpotifyPlaylist(allYoutubeTitles []AllYoutubeItems,spotifyClient *spotify.Client, newSpotifyPlaylist *spotify.FullPlaylist, ytClient *http.Client,newYoutubePlaylist *youtube.Playlist ) error {
	for i, items := range allYoutubeTitles {
		fmt.Println("index: ", i)
		fmt.Println("title: ", items.Title)

		// Search for the video in Spotify
		results, err := spotifyClient.Search(context.Background(), items.Title, spotify.SearchTypeTrack)
		if err != nil {
			fmt.Println("Error searching for video:", err)
			continue
		}

		// Get the first track from the search results
		if len(results.Tracks.Tracks) > 0 {
			firstTrack := results.Tracks.Tracks[0]

			// Add the track to the spotify playlist
			_, err := spotifyClient.AddTracksToPlaylist(context.Background(), newSpotifyPlaylist.ID, firstTrack.ID)
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

	return nil
}

func (r *Resource) addTracksToYoutubePlaylist(allSpotifyTitles []spotify.PlaylistItemTrack,spotifyClient *spotify.Client, newSpotifyPlaylist *spotify.FullPlaylist, ytClient *http.Client,newYoutubePlaylist *youtube.Playlist ) error {
	for i, title := range allSpotifyTitles {
		fmt.Println("index: ", i)
		fmt.Println("title: ", title.Track.Artists[0].Name)
		// Add the track to the newly created playlist
		_, err := spotifyClient.AddTracksToPlaylist(context.Background(), newSpotifyPlaylist.ID, title.Track.ID)
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

	return nil
}