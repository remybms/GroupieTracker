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
	index int
	value int
	flag  bool
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

var artistsDataPaginate artistsPaginate

func paginatehomeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./static/html/Home.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	page_begin := r.FormValue("Button")
	page_Next := r.FormValue("NextButton")
	page_Prev := r.FormValue("PrevButton")
	if len(page_begin) > 0 || len(page_Next) > 0 || len(page_Prev) > 0 {
		if page_begin == "afficher" {
			nbItems := r.FormValue("nb-items")
			First(nbItems)
			artistsDataPaginate.flag = true
			t.Execute(w, artistsDataPaginate)
		} else if page_Next == "Suivant" {
			if artistsDataPaginate.value != 52 && artistsDataPaginate.value != 0 {
				Next()
				t.Execute(w, artistsDataPaginate)
			} else {
				t.Execute(w, artistsData)
			}
			artistsDataPaginate.flag = true
		} else if page_Prev == "Précedent" {
			if artistsDataPaginate.value != 52 && artistsDataPaginate.value != 0 {
				Prev()
				t.Execute(w, artistsDataPaginate)
			} else {
				t.Execute(w, artistsData)
			}
			artistsDataPaginate.flag = false
		}
	} else {
		t.Execute(w, artistsData)
	}
}

func First(nbItems string) {
	artistsDataPaginate.value, _ = strconv.Atoi(nbItems)
	artistsDataPaginate.Array = artistsData.Array[:artistsDataPaginate.value]
	artistsDataPaginate.index = artistsDataPaginate.value
}
func Next() {
	if artistsDataPaginate.flag {
		artistsDataPaginate.index += artistsDataPaginate.value
	} else {
		artistsDataPaginate.index += (artistsDataPaginate.value * 2)
	}
	if (artistsDataPaginate.index) < len(artistsData.Array) {
		artistsDataPaginate.Array = artistsData.Array[(artistsDataPaginate.index - artistsDataPaginate.value):artistsDataPaginate.index]
	} else {
		artistsDataPaginate.index -= artistsDataPaginate.value
		artistsDataPaginate.Array = artistsData.Array[artistsDataPaginate.index:]
	}
}

func Prev() {
	if !artistsDataPaginate.flag {
		artistsDataPaginate.index -= artistsDataPaginate.value
	} else {
		artistsDataPaginate.index -= (artistsDataPaginate.value * 2)
	}
	if (artistsDataPaginate.index) > 0 {
		artistsDataPaginate.Array = artistsData.Array[artistsDataPaginate.index:(artistsDataPaginate.index + artistsDataPaginate.value)]
	} else {
		artistsDataPaginate.index += artistsDataPaginate.value
		artistsDataPaginate.Array = artistsData.Array[:artistsDataPaginate.value]
	}
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

func main() {
	fmt.Println("http://localhost:8080")
	feedData()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", paginatehomeHandler)
	http.HandleFunc("/artist", artistHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/concert", concertHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
