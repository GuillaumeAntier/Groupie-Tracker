package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
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

type Events struct {
	Index []Event 
}

type LocationDates struct {
	Location string
	Dates    []string
}

type LocationCity struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}

var client = &http.Client{
    Timeout: time.Second * 10,  
}

func openAPI(w http.ResponseWriter, r *http.Request, url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		codeErreur(w, r, 500, "Server API is not responding")
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		codeErreur(w, r, 500, "No data to sent")
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

func sortAlphabetically(w http.ResponseWriter,r *http.Request,artists []Artist) []Artist {
	if r.URL.Query().Get("sort") == "A-Z" {
		sort.Slice(artists, func(i, j int) bool {
			return artists[i].Name < artists[j].Name
		})
	} else if r.URL.Query().Get("sort") == "Z-A" {
		sort.Slice(artists, func(i, j int) bool {
			return artists[i].Name > artists[j].Name
		})
	}
	return artists
}

func sortAndFilter(w http.ResponseWriter, r *http.Request, artists []Artist) []Artist {
    research := r.URL.Query().Get("research")
    sortParam := r.URL.Query().Get("sort")

    if research != "" {
        var filteredArtists []Artist
        for _, artist := range artists {
            if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(research)) {
                filteredArtists = append(filteredArtists, artist)
            }
        }
        artists = filteredArtists
    }

    if sortParam == "A-Z" {
        sort.Slice(artists, func(i, j int) bool {
            return artists[i].Name < artists[j].Name
        })
    } else if sortParam == "Z-A" {
        sort.Slice(artists, func(i, j int) bool {
            return artists[i].Name > artists[j].Name
        })
    }

    return artists
}

func sortDatesLocations(datesLocations map[string][]string) []LocationDates {
	var locationDatesSlice []LocationDates
	for location, dates := range datesLocations {
		locationDatesSlice = append(locationDatesSlice, LocationDates{location, dates})
	}
	sort.Slice(locationDatesSlice, func(i, j int) bool {
		layout := "02-01-2006"
		date1, _ := time.Parse(layout, locationDatesSlice[i].Dates[0])
		date2, _ := time.Parse(layout, locationDatesSlice[j].Dates[0])
		return date1.Before(date2)
	})
	return locationDatesSlice
}

func servePageArtist(w http.ResponseWriter, r *http.Request, html string, data []Artist) {
	if r.URL.Path == "/index" {
	} else {
		codeErreur(w, r, 404, "Page not found")
		return
	}

	data = sortAlphabetically(w,r,data)

	page, err := template.ParseFiles("HTML/" + html)
	if err != nil {
		codeErreur(w, r, 500, "Template not found : index.html")
		return
	}
	err = page.Execute(w, data)
	if err != nil {
		codeErreur(w, r, 500, "Internal server error")
		return
	}
}
func servePagePresentation(w http.ResponseWriter, r *http.Request, html string, data []Artist) {
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

func servePagePresentation2(w http.ResponseWriter, r *http.Request, html string, data []Artist) {
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

func servePageResult(w http.ResponseWriter, r *http.Request, html string, data []Artist) {
	if r.URL.Path != "/result" {
		codeErreur(w, r, 404, "Page not found")
		return
	}

	data = sortAndFilter(w,r,data)

	page, err := template.ParseFiles("HTML/" + html)
	if err != nil {
		codeErreur(w, r, 500, "Template not found : result.html")
		return
	}
	err = page.Execute(w, data)
	if err != nil {
		codeErreur(w, r, 500, "Internal server error")
		return
	}
}

func servePageEvent(w http.ResponseWriter, r *http.Request, html string, dataArtist []Artist ,data []LocationDates, coordinatesMap map[string][]LocationCity) {
	if r.URL.Path != "/event" {
		codeErreur(w, r, 404, "Page not found")
		return
	}
	page, err := template.ParseFiles("HTML/" + html)
	if err != nil {
		codeErreur(w, r, 500, "Template not found : event.html")
		return
	}
	err = page.Execute(w, struct {
		Artists        []Artist
		DatesLocations []LocationDates
		Coordinates    map[string][]LocationCity
	}{	
		Artists:        dataArtist,
		DatesLocations: data,
		Coordinates:    coordinatesMap,
	})
	if err != nil {
		codeErreur(w, r, 500, "Internal server error")
		return
	}
}



func Home(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyart, err := openAPI(w,r,url + "artists")
	var artist []Artist
	err = json.Unmarshal(bodyart, &artist)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	servePagePresentation(w, r, "presentation.html", artist)

}

func Home2(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyart, err := openAPI(w,r,url + "artists")
	var artist []Artist
	err = json.Unmarshal(bodyart, &artist)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	servePagePresentation2(w, r, "presentation2.html", artist)
}

func HandlerMain(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyart, err := openAPI(w, r, url+"artists")
	var artist []Artist
	err = json.Unmarshal(bodyart, &artist)
	if err != nil {
		codeErreur(w, r, 500, "Internal server error")
		return
	}
	servePageArtist(w, r, "index.html", artist)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	research := r.URL.Query().Get("research")
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyart, err := openAPI(w, r, url+"artists")
	var artist []Artist
	err = json.Unmarshal(bodyart, &artist)
	if err != nil {
		codeErreur(w, r, 500, "Internal server error")
		return
	}
	artist = filterArtistsByLetter(artist, research)
	servePageResult(w, r, "result.html", artist)
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyevent, err := openAPI(w, r, url+"relation/"+id)
	var event Event
	err = json.Unmarshal(bodyevent, &event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bodyart, err := openAPI(w, r, url+"artists")
	var art []Artist
	err = json.Unmarshal(bodyart, &art)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var artist []Artist
	for _, a := range art {
		id,err := strconv.Atoi(id)
		if err != nil {
			codeErreur(w, r, 500, "Internal server error")
			return
		}
		if a.ID == id {
			artist = append(artist, a)
		}
	}
	locationDateSlice := sortDatesLocations(event.DatesLocations)
	cordinatesMap := make(map[string][]LocationCity)

	for location, _ := range event.DatesLocations {
		latitude, longitude, err := getCoordinates(location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cordinatesMap[location] = []LocationCity{{latitude, longitude}}
	}
	servePageEvent(w, r, "event.html", artist ,locationDateSlice, cordinatesMap)
}


func getCoordinates(location string) (string, string, error) {
	params := url.Values{}
	location = strings.Replace(location, "-", ",", -1)
	params.Set("q", location)
	params.Set("format", "json")

	resp, err := client.Get("https://nominatim.openstreetmap.org/search?format=json&q=" + url.QueryEscape(location))
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

func codeErreur(w http.ResponseWriter, r *http.Request, status int, message string) {

	colorRed := "\033[31m"

	if status == 404 {
		http.Error(w, "404 not found", http.StatusNotFound)
		fmt.Println(string(colorRed), "[SERVER_ALERT] - 404 : File not found, or missing...", message)
		if status == 400 {
			http.Error(w, "400 Bad request", http.StatusBadRequest)
			fmt.Println(string(colorRed), "[SERVER_ALERT] - 400 : Bad request", message)
		}
		if status == 500 {
			http.Error(w, "500 Internal server error", http.StatusInternalServerError)
			fmt.Println(string(colorRed), "[SERVER_ALERT] - 500 : Internal server error", message)
		}

	}
}

func main() {
	fmt.Println(string("\033[34m"), "[SERVER_INFO] : Starting local Server...")
	
	http.Handle("/JS/", http.StripPrefix("/JS/", http.FileServer(http.Dir("JS"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/presentation", Home)
	http.HandleFunc("/presentation2", Home2)
	http.HandleFunc("/", HandlerMain)
	http.HandleFunc("/index", HandlerMain)
	http.HandleFunc("/result", searchHandler)
	http.HandleFunc("/event", eventHandler)

	fmt.Println(string("\033[32m"), "[SERVER_READY] : on http://localhost:8080/index ✅ ")
	fmt.Println(string("\033[33m"), "[SERVER_INFO] : To stop the program : Ctrl + c")
	log.Fatal(http.ListenAndServe(":8080", nil))

}