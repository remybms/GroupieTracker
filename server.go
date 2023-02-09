package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
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

var artistsData []artists

func artist() {

	url := "https://groupietrackers.herokuapp.com/api/artists"
	req, _ := http.NewRequest("GET", url, nil)
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	//fmt.Println(string(body))
	err := json.Unmarshal([]byte(body), &artistsData)
	defer res.Body.Close()
	if err != nil {
		fmt.Println("Error :", err)
		return
	}
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./static/artist.html")
	for index := range artistsData {
		t.Execute(w, artistsData[index])
	}
}

func main() {
	fmt.Println("http://localhost:8080")
	artist()
	fmt.Println(artistsData)
	http.HandleFunc("/", artistHandler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
