package main

import (

	//"os"
	//"bufio"
	"html/template"
	"io/ioutil"
	"strings"

	//"encoding/json"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var response string
var respart []string
var resploc []string
var respdat []string
var resprel []string
var err []string
var name string
var artist Artist

type Artist struct{
	ID int `json:"id"`
	Image string `json:"image"`
	Name string `json:"name"`
	Members []string `json:"members"`
	CreationDate int `json:"creationDate"`
	FirstAlbum string `json:"firstAlbum"`
}

type Event struct{
	DatesLocations map[string][]string
}

type Relation struct {
	Artist []Artist 
	Event []Event
}
func openAPI(url string) ([]byte, error)  {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return nil ,err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture:", err)
		return nil ,err
	}
	return body, nil
}

func filterArtistsByLetter(artists []Artist, letter string) []Artist {
	filteredArtists:= make([]Artist, 0)
	for _, artist := range artists {
		if strings.HasPrefix(strings.ToLower(artist.Name), strings.ToLower(letter)) {
				filteredArtists = append(filteredArtists, artist)
		}
	}
	return filteredArtists
}

func servePageArtist (w http.ResponseWriter, r *http.Request, html string, data []Artist) {
	page,err := template.ParseFiles("HTML/"+html)
	if err != nil {
		fmt.Println(err)
	}
	err = page.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func servePage (w http.ResponseWriter, r *http.Request, html string, data Event) {
	page,err := template.ParseFiles("HTML/"+html)
	if err != nil {
		fmt.Println(err)
	}
	err = page.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func HandlerMain(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyart, err := openAPI(url +"artists")
	var artist []Artist
	err = json.Unmarshal(bodyart, &artist)
	if err != nil {
		fmt.Println("Erreur lors de la lecture:", err)
		return
	}
	servePageArtist(w, r, "index.html", artist)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	research := r.URL.Query().Get("research")
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyart, err := openAPI(url +"artists")
	var artist []Artist
	err = json.Unmarshal(bodyart, &artist)
	if err != nil {
		fmt.Println("Erreur lors de la lecture:", err)
		return
	}
	artist = filterArtistsByLetter(artist, research)
	servePageArtist(w, r, "result.html", artist)

}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyevent, err := openAPI(url +"relation/"+id)
	var event Event
	err = json.Unmarshal(bodyevent, &event)
	if err != nil {
		fmt.Println("Erreur lors de la lecture:", err)
		return
	}
	servePage(w, r, "event.html", event)
}

func main() {
	
	http.HandleFunc("/", HandlerMain)
	http.HandleFunc("/index", HandlerMain)
	http.HandleFunc("/result", searchHandler)
	http.HandleFunc("/event", eventHandler)

	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}