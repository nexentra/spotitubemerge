package spotifyplaylist

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	spotify "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func RegisterHandlers(r *echo.Group, authenticator *spotifyauth.Authenticator, redirectURI string, state string, errorLog *log.Logger, infoLog *log.Logger, env map[string]string) {
	res := Resource{
		Authenticator: authenticator,
		RedirectURI:   redirectURI,
		State:         state,
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		Env:           env,
	}

	r.GET("/auth/spotify", res.loginSpotify)
	r.POST("/auth/spotify/callback", res.callbackSpotify)
	r.GET("/spotify-playlist", res.getSpotifyPlaylist)
	r.GET("/spotify-items", res.getSpotifyItems)
	r.GET("/search-spotify-items", res.searchSpotifyItems)
}

type Resource struct {
	Authenticator *spotifyauth.Authenticator
	Client        *spotify.Client
	RedirectURI   string
	State         string
	UserId        string
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	Env           map[string]string
}

func (r Resource) loginSpotify(c echo.Context) error {
	authURL := r.Authenticator.AuthURL(r.State)
	fmt.Println("Auth URL: ", authURL)
	return c.JSON(http.StatusOK, echo.Map{
		"authUrl": authURL,
	})
}

type SpotifyCode struct {
	Code string `json:"code"`
}

func (r Resource) callbackSpotify(c echo.Context) error {
	spotifyCode := SpotifyCode{}
	if err := c.Bind(&spotifyCode); err != nil {
		r.ErrorLog.Printf("Failed to bind code: %v", err)
	}
	fmt.Println("Code: ", spotifyCode.Code)

	// state := r.State
	tok, err := r.Authenticator.Exchange(c.Request().Context(), spotifyCode.Code)
	if err != nil {
		r.ErrorLog.Print(err)
	}
	fmt.Println("Token: ", tok)

	// if st := c.FormValue("state"); st != state {
	// 	r.ErrorLog.Printf("State mismatch: %s != %s\n", st, state)
	// }

	client := spotify.New(r.Authenticator.Client(c.Request().Context(), tok))
	fmt.Println("Client: ", client)
	r.Client = client

	user, err := r.Client.CurrentUser(c.Request().Context())
	if err != nil {
		r.ErrorLog.Print(err)
	}
	r.UserId = user.ID
	fmt.Println("You are logged in as:", user.ID)
	return c.JSON(http.StatusOK, echo.Map{
		"token": tok,
	})
}

func (r Resource) getSpotifyPlaylist(c echo.Context) error {
	var authHeaderType *oauth2.Token
	authHeader := c.Request().Header.Get("AuthorizationSpotify")
	json.Unmarshal([]byte(authHeader), &authHeaderType)
	client := spotify.New(r.Authenticator.Client(c.Request().Context(), authHeaderType))

	user, err := client.CurrentUser(c.Request().Context())
	if err != nil {
		r.ErrorLog.Print(err)
	}

	var allPlaylists []spotify.SimplePlaylist
	for {
		playlists, err := client.GetPlaylistsForUser(c.Request().Context(), user.ID, spotify.Limit(50), spotify.Offset(len(allPlaylists)))
		if err != nil {
			r.ErrorLog.Print(err)
		}
		if err != nil {
			r.ErrorLog.Print(err)
			break
		}

		allPlaylists = append(allPlaylists, playlists.Playlists...)

		// Check if there are more playlists to fetch
		if playlists.Next == "" {
			fmt.Println("break")
			break
		}
	}

	return c.JSON(http.StatusOK, echo.Map{
		"playlists": allPlaylists,
	})
}

func (r Resource) getSpotifyItems(c echo.Context) error {
	var authHeaderType *oauth2.Token
	authHeader := c.Request().Header.Get("AuthorizationSpotify")
	json.Unmarshal([]byte(authHeader), &authHeaderType)
	strings := c.QueryParam("strings")
	fmt.Println("strings: ", strings)

	if strings == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": "strings is empty",
		})
	}

	client := spotify.New(r.Authenticator.Client(c.Request().Context(), authHeaderType))

	// user, err := client.CurrentUser(c.Request().Context())
	// if err != nil {
	// 	r.ErrorLog.Print(err)
	// }

	playlist, err := client.GetPlaylistItems(c.Request().Context(), spotify.ID(strings))
	if err != nil {
		r.ErrorLog.Print(err)
	}
	for _, playlist := range playlist.Items {
		fmt.Println("  ", playlist.Track)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"items": playlist,
	})
}

func (r Resource) searchSpotifyItems(c echo.Context) error {
	var authHeaderType *oauth2.Token
	authHeader := c.Request().Header.Get("AuthorizationSpotify")
	json.Unmarshal([]byte(authHeader), &authHeaderType)

	strings := c.QueryParam("strings")
	fmt.Println("strings: ", strings)

	// searching tracks with given name
	client := spotify.New(r.Authenticator.Client(c.Request().Context(), authHeaderType))
	results, err := client.Search(c.Request().Context(), strings, spotify.SearchTypeTrack) //spotify.SearchTypePlaylist|spotify.SearchTypeAlbum
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("results:", results)

	// handle songs results
	if results.Tracks != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"items": results.Tracks.Tracks,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"items": []string{},
	})
}
