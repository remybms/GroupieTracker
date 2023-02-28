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
type artistsArray struct {
	*artists
	Array []artists
	Valid []artists
	Flag  bool
}

var artistsData artistsArray

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

func homeHandler(w http.ResponseWriter, r *http.Request) {
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
	t, err := template.ParseFiles("./static/html/Research.html")
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

func artistHandler(w http.ResponseWriter, r *http.Request) {
	indexString := r.FormValue("card")
	indexStringSelect := r.FormValue("languages")
	index := 0
	if len(indexString) > 0 {
		index, _ = strconv.Atoi(indexString)
	} else {
		index, _ = strconv.Atoi(indexStringSelect)
	}
	t, err := template.ParseFiles("./static/html/Artist.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	t.Execute(w, artistsData.Array[index-1])
}

func main() {
	fmt.Println("http://localhost:8080")
	Artists()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/artist", artistHandler)
	http.HandleFunc("/search", searchHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
