package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const (
	clientID     = "b060eb29f3d44f388961e66fd7b55fa4"
	clientSecret = "4bd35b387c0446aab01c3c4961c1439d"
	redirectURI  = "http://localhost:8080/callback"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type Artist struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Genres     []string `json:"genres"`
	Followers  int      `json:"followers"`
	Popularity int      `json:"popularity"`
}

type Track struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	PreviewURL string `json:"preview_url"`
}

type SearchResponse struct {
	Artists []Artist `json:"artists"`
	Tracks  []Track  `json:"tracks"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/html/testspotify.html")
	if err != nil {
		return
	}
	code := r.URL.Query().Get("code")
	token, _ := getAccessToken(code)
	authURL := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=%s&redirect_uri=http://localhost:8080/callback&scope=user-read-private+user-read-email+playlist-read-private", clientID, token)
	http.Redirect(w, r, authURL, http.StatusFound)
	t.Execute(w, http.StatusFound)
}

func getAccessToken(code string) (*Token, error) {
	endpoint := fmt.Sprintf("https://accounts.spotify.com/api/token?grant_type=authorization_code&code=%s&redirect_uri=%s", code, redirectURI)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(clientID, clientSecret)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
func callbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Failed to get authorization code", http.StatusBadRequest)
		return
	}

	token, err := getAccessToken(code)
	if err != nil {
		http.Error(w, "Failed to get access token", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Access token: %s\n", token.AccessToken)

	artist := r.URL.Query().Get("artist")
	if artist == "" {
		fmt.Fprintln(w, "Please provide an artist name")
		return
	}

	searchResult, err := searchArtist(token.AccessToken, artist)
	if err != nil {
		http.Error(w, "Failed to search for artist", http.StatusInternalServerError)
		return
	}

	if len(searchResult.Artists) == 0 {
		fmt.Fprintf(w, "No artists found for \"%s\"\n", artist)
		return
	}

	topTracks, err := getArtistTopTracks(token.AccessToken, searchResult.Artists[0].ID)
	if err != nil {
		http.Error(w, "Failed to get artist's top tracks", http.StatusInternalServerError)
		return
	}

	if len(topTracks) == 0 {
		fmt.Fprintf(w, "No top tracks found for \"%s\"\n", searchResult.Artists[0].Name)
		return
	}

	fmt.Fprintf(w, "Playing preview of \"%s\" by \"%s\"\n", topTracks[0].Name, searchResult.Artists[0].Name)
	fmt.Fprintf(w, "Preview URL: %s\n", topTracks[0])
}

func getArtistTopTracks(accessToken string, artistID string) ([]Track, error) {
	endpoint := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/top-tracks?market=US", artistID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var topTracksResponse struct {
		Tracks []Track `json:"tracks"`
	}

	err = json.Unmarshal(body, &topTracksResponse)
	if err != nil {
		return nil, err
	}

	return topTracksResponse.Tracks, nil
}
func searchArtist(accessToken string, artistName string) (*SearchResponse, error) {
	endpoint := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&type=artist", url.QueryEscape(artistName))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var searchResponse SearchResponse
	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return nil, err
	}

	return &searchResponse, nil
}
func main() {
	fmt.Println("http://localhost:8080")
	http.HandleFunc("/", handler)
	http.HandleFunc("/callback", callbackHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
