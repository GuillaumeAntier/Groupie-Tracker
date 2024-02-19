package main

import (

	"html/template"
	"io/ioutil"
	"strings"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sort"
	"time"
)

type Artist struct{
	ID int `json:"id"`
	Image string `json:"image"`
	Name string `json:"name"`
	Members []string `json:"members"`
	CreationDate int `json:"creationDate"`
	FirstAlbum string `json:"firstAlbum"`
}

type Event struct{
	DatesLocations map[string][]string `json:"datesLocations"`
}

type LocationDates struct{
	Location string
	Dates []string
}

type LocationCity struct {
    Latitude  string `json:"lat"`
    Longitude string `json:"lon"`
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
	page,err := template.ParseFiles("HTML/"+html)
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}
	err = page.Execute(w, data)
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}
}	

func servePageEvent(w http.ResponseWriter, r *http.Request, html string, data []LocationDates, coordinatesMap map[string][]LocationCity) {
	page,err := template.ParseFiles("HTML/"+html)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = page.Execute(w, struct {
    	DatesLocations []LocationDates
    	Coordinates map[string][]LocationCity
	}{
    	DatesLocations : data,
    	Coordinates : coordinatesMap,
    
	})
	if err != nil {
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return 
	}
}


func HandlerMain(w http.ResponseWriter, r *http.Request) {
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyart, err := openAPI(url +"artists")
	var artist []Artist
	err = json.Unmarshal(bodyart, &artist)
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
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
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	locationDateSlice := sortDatesLocations(event.DatesLocations)	
	cordinatesMap := make(map[string][]LocationCity)

	for location, _ := range event.DatesLocations {
		latitude, longitude, err := getCoordinates(location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cordinatesMap[location] = []LocationCity{{latitude,longitude}}
	}
	servePageEvent(w, r, "event.html", locationDateSlice, cordinatesMap)
}

func getCoordinates(location string) (string, string, error) {
    baseURL := "https://nominatim.openstreetmap.org/search"
    params := url.Values{}
	location = strings.Replace(location, "-", ",", -1)
    params.Set("q", location)
    params.Set("format", "json")

    resp, err := http.Get(baseURL + "?" + params.Encode())
    if err != nil {
        return "","", err
    }
    defer resp.Body.Close()

	var results []LocationCity
    err = json.NewDecoder(resp.Body).Decode(&results)
    if err != nil {
        return "","", err
    }

    if len(results) == 0 {
        return "","", fmt.Errorf("No results found for location: %s", location)
    }

    return results[0].Latitude, results[0].Longitude, nil
}

func main() {
	http.Handle("/JS/", http.StripPrefix("/JS/", http.FileServer(http.Dir("JS"))))
	http.HandleFunc("/", HandlerMain)
	http.HandleFunc("/index", HandlerMain)
	http.HandleFunc("/result", searchHandler)
	http.HandleFunc("/event", eventHandler)


	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}