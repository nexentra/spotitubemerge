package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	ctx    context.Context
	config *oauth2.Config
	token  *oauth2.Token
)

func main() {
	// Initialize the web app
	initWebApp()

	http.HandleFunc("/", handleHome)
	http.HandleFunc("/auth", handleAuth)
	http.HandleFunc("/callback", handleCallback)
	http.HandleFunc("/playlists", handlePlaylists)

	log.Println("Starting server on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initWebApp() {
	ctx = context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err = google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the YouTube Playlists Web App!")
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOnline)
	http.Redirect(w, r, authURL, http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	var code = r.URL.Query().Get("code")
	tok, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	// Save the token in the variable instead of caching
	token = tok

	fmt.Println("Token: ", token.AccessToken)
	getPlaylists()
	// http.Redirect(w, r, "/playlists", http.StatusFound)
}


func getPlaylists() {
	url := "https://www.googleapis.com/youtube/v3/playlists?part=snippet%2CcontentDetails&maxResults="+"25"+"&mine=true&key="+ config.ClientID+ "&access_token=" + token.AccessToken

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %v", err)
	}

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	fmt.Println(string(body))
}


func handlePlaylists(w http.ResponseWriter, r *http.Request) {
	if token == nil {
		http.Redirect(w, r, "/auth", http.StatusFound)
		return
	}

	// fmt.Println("Token: ", token)

	// client := config.Client(ctx, token)
	// if client == nil {
	// 	log.Fatal("Unable to create OAuth client")
	// }
	// fmt.Println("client: ", client)
	ctx = context.Background()
	service, err := youtube.NewService(
		ctx,
		option.WithHTTPClient(&http.Client{
			Timeout: 30 * time.Second,
			Transport: &oauth2.Transport{
				Base:   http.DefaultTransport,
				Source: config.TokenSource(ctx, token),
			},
		}),
	)
	handleError(err, "Error creating YouTube client")
	fmt.Println("service: ", service)

	part := []string{"snippet,contentDetails,statistics"}
	mine := false
	call := service.Playlists.List(part)
	call = call.Mine(mine)
	response, err := call.Do()
	handleError(err, "")

	fmt.Fprintln(w, "YouTube Playlists:")
	for _, playlist := range response.Items {
		fmt.Fprintln(w, "  "+playlist.Snippet.Title)
	}
}

func handleError(err error, message string) {
	if message != "" && err != nil {
		log.Fatalf(message+": %v", err)
	}
}