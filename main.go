package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
}

type Event struct {
	DatesLocations map[string][]string `json:"datesLocations"`
}

type Relation struct {
	Artist []Artist
	Event  []Event
}

type LocationCity struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}

func openAPI(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture:", err)
		return nil, err
	}
	return body, nil
}

func filterArtistsByLetter(artists []Artist, letter string) []Artist {
	filteredArtists := make([]Artist, 0)
	for _, artist := range artists {
		if strings.HasPrefix(strings.ToLower(artist.Name), strings.ToLower(letter)) {
			filteredArtists = append(filteredArtists, artist)
		}
	}
	return filteredArtists
}

func servePageArtist(w http.ResponseWriter, r *http.Request, html string, data []Artist) {
	page, err := template.ParseFiles("html/" + html)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	err = page.Execute(w, data)
	if err != nil {
		http.NotFound(w, r)
		return
	}
}

func servePage(w http.ResponseWriter, r *http.Request, html string, data Event) {
	page, err := template.ParseFiles("html/" + html)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	err = page.Execute(w, data)
	if err != nil {
		http.NotFound(w, r)
		return
	}
}

func HandlerMain(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyart, err := openAPI(url + "artists")
	var artist []Artist
	err = json.Unmarshal(bodyart, &artist)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	servePageArtist(w, r, "index.html", artist)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	research := r.URL.Query().Get("research")
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyart, err := openAPI(url + "artists")
	var artist []Artist
	err = json.Unmarshal(bodyart, &artist)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	artist = filterArtistsByLetter(artist, research)
	servePageArtist(w, r, "result.html", artist)
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyevent, err := openAPI(url + "relation/" + id)
	var event Event
	err = json.Unmarshal(bodyevent, &event)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	for _, locations := range event.DatesLocations {
		latitude, longitude, err := getCoordinates(locations[0])
		if err != nil {
			http.NotFound(w, r)
			return
		}
		fmt.Println("Latitude:", latitude, "Longitude:", longitude)
	}
	servePage(w, r, "event.html", event)
}

func getCoordinates(location string) (string, string, error) {
	baseURL := "https://nominatim.openstreetmap.org/search"
	params := url.Values{}
	params.Set("q", location)
	params.Set("format", "json")

	resp, err := http.Get(baseURL + "?" + params.Encode())
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var results []LocationCity
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return "", "", err
	}

	if len(results) == 0 {
		return "", "", fmt.Errorf("No results found for location: %s", location)
	}

	return results[0].Latitude, results[0].Longitude, nil
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", HandlerMain)
	http.HandleFunc("/index.html", HandlerMain)
	http.HandleFunc("/result.html", searchHandler)
	http.HandleFunc("/event.html", eventHandler)

	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
