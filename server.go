package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
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

type rangeRelation struct {
	Name     string
	Location []string
	Dates    [][]string
}

type ExtractRelation struct {
	Index []relation `json:"index"`
}

type ForBingAPI struct {
	ResourceSets []struct {
		Resources []struct {
			Point struct {
				Coordinates []float64
			}
		}
	}
}

type coordinates struct {
	Name      string
	Latitude  float64
	Longitude float64
}

type artistsArray struct {
	*artists
	Array []artists
	Valid []artists
	Flag  bool
}

var Maps ForBingAPI
var coordinatesMap coordinates
var artistsData artistsArray
var concertsData ExtractRelation

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

func Map(location string) {
	url := fmt.Sprintf("https://dev.virtualearth.net/REST/v1/Locations?q=%s&key=%s", location, "AlBlNdfGSHdDQO7QSc9vamIHHUD6c0VArZIZ9i3l-F9J4whlFM9Fz3ZMxE1t_lMh")
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	err := json.Unmarshal([]byte(body), &Maps)
	if err != nil {
		fmt.Println("Error :", err)
		return
	}
	defer res.Body.Close()
}

func Relation() {

	url := "https://groupietrackers.herokuapp.com/api/relation"
	req, errorRequest := http.NewRequest("GET", url, nil)
	if errorRequest != nil {
		log.Fatal(errorRequest)
	}
	res, errorServ := http.DefaultClient.Do(req)
	if errorServ != nil {
		log.Fatal(errorServ)
	}
	body, errorFich := ioutil.ReadAll(res.Body)
	if errorFich != nil {
		log.Fatal(errorFich)
	}
	//fmt.Println(string(body))
	err := json.Unmarshal([]byte(body), &concertsData)

	if err != nil {
		fmt.Println("Error :", err)
		return
	}
	defer res.Body.Close()
}

func feedData() {
	Artists()
	Relation()
}

func defineOrder(order string) {
	if order == "creationDate" {
		sort.Slice(artistsData.Array[:], func(i, j int) bool {
			return artistsData.Array[i].CreationDate < artistsData.Array[j].CreationDate
		})
	} else if order == "id" {
		sort.Slice(artistsData.Array[:], func(i, j int) bool {
			return artistsData.Array[i].Id < artistsData.Array[j].Id
		})
	} else if order == "a-z" {
		sort.Slice(artistsData.Array[:], func(i, j int) bool {
			return artistsData.Array[i].Name < artistsData.Array[j].Name
		})
	} else if order == "z-a" {
		sort.Slice(artistsData.Array[:], func(i, j int) bool {
			return artistsData.Array[i].Name > artistsData.Array[j].Name
		})
	} else if order == "reverseCreationDate" {
		sort.Slice(artistsData.Array[:], func(i, j int) bool {
			return artistsData.Array[i].CreationDate > artistsData.Array[j].CreationDate
		})
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tri := r.FormValue("tri")
	defineOrder(tri)
	t, err := template.ParseFiles("./static/html/Home.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Execute(w, artistsData)

}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("tu es dans la boucle")
	artistsData.artists = &artistsData.Array[0]
	indexString := r.FormValue("research")
	artistsData.Valid = []artists{}
	artistsData.Flag = false
	t, err := template.ParseFiles("./static/html/research.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, value := range artistsData.Array {
		val := strings.ToLower(value.Name)
		str := strings.ToLower(indexString)
		if strings.Contains(val, str) {
			artistsData.Valid = append(artistsData.Valid, value)
			artistsData.Flag = true
		}
	}
	t.Execute(w, artistsData)
}

var index int

func concertHandler(w http.ResponseWriter, r *http.Request) {
	var concertsDatesLocations rangeRelation
	indexString := r.FormValue("dates")
	if indexString != "" {
		index, _ = strconv.Atoi(indexString)
	}
	t, err := template.ParseFiles("./static/html/concert.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	concertsDatesLocations.Name = string(artistsData.Array[index-1].Name)
	for location, dates := range concertsData.Index[index-1].DatesLocations {
		concertsDatesLocations.Location = append(concertsDatesLocations.Location, location)
		concertsDatesLocations.Dates = append(concertsDatesLocations.Dates, dates)
	}
	t.Execute(w, concertsDatesLocations)

}

func mapHandler(w http.ResponseWriter, r *http.Request) {
	location := r.FormValue("location")
	Map(location)
	t, err := template.ParseFiles("./static/html/Map.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	coordinatesMap.Name = string(artistsData.Array[index-1].Name)
	coordinatesMap.Latitude = Maps.ResourceSets[0].Resources[0].Point.Coordinates[0]
	coordinatesMap.Longitude = Maps.ResourceSets[0].Resources[0].Point.Coordinates[1]
	t.Execute(w, coordinatesMap)
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	indexString := r.FormValue("card")
	if indexString != "" {
		index, _ = strconv.Atoi(indexString)
	}
	sort.Slice(artistsData.Array[:], func(i, j int) bool {
		return artistsData.Array[i].Id < artistsData.Array[j].Id
	})
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
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist", artistHandler)
	http.HandleFunc("/search", searchHandler)
	http.HandleFunc("/concert", concertHandler)
	http.HandleFunc("/map", mapHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
