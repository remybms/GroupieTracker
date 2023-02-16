package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
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
	DatesLocations relationDates
}

type locations struct {
	Id        int
	Locations []string
	Dates     string
}

type dates struct {
	Id    int
	Dates []string
}

type relationDates struct {
	DatesToLocations []string
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
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/html/Home.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Execute(w, artistsData)

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

func main() {
	fmt.Println("http://localhost:8080")
	feedData()
	fmt.Println(concertsData)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist", artistHandler)
	http.HandleFunc("/search", searchHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
