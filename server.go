package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type artists struct {
	Id           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Location     string   `json:"location"`
	ConcertDates string   `json:"concertDates"`
	Relations    string   `json:"relations"`
}

type relation struct {
	Id             int
	DatesLocations map[string][]string
}

type locations struct {
	Id        int
	Locations []string
	Dates     string
}

type rangeRelation struct {
	Location []string
	Dates    [][]string
}

type dates struct {
	Id    int
	Dates []string
}

type ExtractDate struct {
	Index []dates `json:"index"`
}

type ExtractLocation struct {
	Index []locations `json:"index"`
}

type ExtractRelation struct {
	Index []relation `json:"index"`
}

type artistsArray struct {
	Array []artists
}

type artistsPaginate struct {
	Array []artists
}

type concerts struct {
	Relation  ExtractRelation
	Locations ExtractLocation
	Dates     ExtractDate
}

var artistsData artistsArray
var concertsData concerts

func Artists() {

	url := "https://groupietrackers.herokuapp.com/api/artists"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	err := json.Unmarshal([]byte(body), &artistsData.Array)

	if err != nil {
		fmt.Println("Error :", err)
		return
	}
	defer res.Body.Close()
}

func Relation() {

	url := "https://groupietrackers.herokuapp.com/api/relation"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	err := json.Unmarshal([]byte(body), &concertsData.Relation)

	if err != nil {
		fmt.Println("Error :", err)
		return
	}
	defer res.Body.Close()
}

func Locations() {

	url := "https://groupietrackers.herokuapp.com/api/locations"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	err := json.Unmarshal([]byte(body), &concertsData.Locations)

	if err != nil {
		fmt.Println("Error :", err)
		return
	}
	defer res.Body.Close()
}

func Dates() {

	url := "https://groupietrackers.herokuapp.com/api/dates"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	err := json.Unmarshal([]byte(body), &concertsData.Dates)

	if err != nil {
		fmt.Println("Error :", err)
		return
	}
	defer res.Body.Close()
}

func feedData() {
	Artists()
	Relation()
	Locations()
	Dates()
	// RunSpotify()
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var indexString string
	var indexRange int = 0
	var artistsDataPaginate artistsPaginate
	t, err := template.ParseFiles("./static/html/Home.html")
	nbItems := r.FormValue("nb-items")
	if nbItems != "" {
		indexString = nbItems
	}
	index, _ := strconv.Atoi(indexString)
	if indexRange < index {
		indexRange = 0
	} else {
		indexRange -= index
	}
	if indexString == "" {
		artistsDataPaginate.Array = artistsData.Array
	} else {
		for nbItem := 0; nbItem < index; nbItem++ {
			if indexRange < len(artistsData.Array) {

				artistsDataPaginate.Array = append(artistsDataPaginate.Array, artistsData.Array[indexRange])
				indexRange++
			}
		}
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Execute(w, artistsDataPaginate)
}

var currentIndex int = 0

func handleNextButton(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/html/Home.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	NextButton := r.FormValue("NextButton")
	if NextButton != "" {
		currentIndex++
		if currentIndex >= len(artistsData.Array) {
			currentIndex = 0
		}
	}

	artistsDataPaginate := artistsPaginate{Array: artistsData.Array[currentIndex : currentIndex+10]}
	t.Execute(w, artistsDataPaginate)
}

func handlePrevButton(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/html/Home.html")
	if err != nil {
		fmt.Println(err)
		return
	}

	PreviousButton := r.FormValue("PreviousButton")
	if PreviousButton != "" {
		currentIndex--
		if currentIndex < 0 {
			currentIndex = len(artistsData.Array) - 1
		}
	}

	artistsDataPaginate := artistsPaginate{Array: artistsData.Array[currentIndex : currentIndex+10]}
	t.Execute(w, artistsDataPaginate)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	indexString := r.FormValue("research")
	fmt.Println(indexString)
	t, err := template.ParseFiles("./static/html/Artist.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	for ind, value := range artistsData.Array {
		if strings.ToLower(value.Name) == strings.ToLower(indexString) {
			t.Execute(w, artistsData.Array[ind])
		}
	}
}

func concertHandler(w http.ResponseWriter, r *http.Request) {
	var concertsDatesLocations rangeRelation
	indexString := r.FormValue("dates")
	index, _ := strconv.Atoi(indexString)
	t, err := template.ParseFiles("./static/html/concert.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	for location, dates := range concertsData.Relation.Index[index-1].DatesLocations {
		concertsDatesLocations.Location = append(concertsDatesLocations.Location, location)
		concertsDatesLocations.Dates = append(concertsDatesLocations.Dates, dates)
	}
	t.Execute(w, concertsDatesLocations)

}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	indexString := r.FormValue("card")
	index, _ := strconv.Atoi(indexString)
	t, err := template.ParseFiles("./static/html/Artist.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Execute(w, artistsData.Array[index-1])

}

// API Spotify//
/*
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

func spotifyHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./testspotify.html")

	authURL := fmt.Sprintf("https://accounts.spotify.com/authorize?client_id=%s&response_type=code&redirect_uri=%s&scope=user-read-private+user-read-email+playlist-read-private", clientID, redirectURI)
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

	var searchResult SearchResponse
	err = json.Unmarshal(body, &searchResult)
	if err != nil {
		return nil, err
	}

	return &searchResult, nil
}
*/
//API Spotify//

/*func main() {
	fmt.Println("http://localhost:8080")
	feedData()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/next", handleNextButton)
	http.HandleFunc("/prev", handlePrevButton)
	http.HandleFunc("/artist", artistHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/concert", concertHandler)
	//http.HandleFunc("/testspotify", spotifyHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}*/
