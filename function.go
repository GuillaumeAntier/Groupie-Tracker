package GroupieTracker

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sort"
	"time"
	"strconv"
	"net/url"
)

var client = &http.Client{
    Timeout: time.Second * 10,  
}

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

func OpenAPI(w http.ResponseWriter, r *http.Request, url string) ([]byte, error) {
	/// OpenAPI is a function that gets the data from the API and returns the data and an error.
	/// It takes the response writer, the request and the url as parameters.
	/// It returns the data and an error.

	res, err := http.Get(url) // Get the data from the API
	if err != nil {
		CodeErreur(w, r, 500, "Server API is not responding")
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body) // Read the data from the API
	if err != nil {
		CodeErreur(w, r, 500, "No data to sent")
		return nil, err
	}
	return body, nil
}

func FilterArtistsByLetter(artists []Artist, letter string) []Artist {
	/// FilterArtistsByLetter is a function that filters the artists by the letter.
	/// It takes the artists and the letter as parameters.
	/// It returns the filtered artists.

	filteredArtists := make([]Artist, 0) // Create a slice to store the filtered artists
	for _, artist := range artists { // Loop through the artists
		if strings.HasPrefix(strings.ToLower(artist.Name), strings.ToLower(letter)) { // Check if the artist name starts with the letter
			filteredArtists = append(filteredArtists, artist) // Add the artist to the filtered artists
		}
	}
	return filteredArtists
}

func sortAlphabetically(w http.ResponseWriter,r *http.Request,artists []Artist) []Artist {
	/// sortAlphabetically is a function that sorts the artists alphabetically.
	/// It takes the response writer, the request and the artists as parameters.
	/// It returns the sorted artists.

	if r.URL.Query().Get("sort") == "A-Z" { // Check if the sort parameter is A-Z
		sort.Slice(artists, func(i, j int) bool {
			return artists[i].Name < artists[j].Name
		})
	} else if r.URL.Query().Get("sort") == "Z-A" { // Check if the sort parameter is Z-A
		sort.Slice(artists, func(i, j int) bool { 
			return artists[i].Name > artists[j].Name
		})
	}
	return artists
}

func SortAndFilter(w http.ResponseWriter, r *http.Request, artists []Artist) []Artist {
	/// SortAndFilter is a function that sorts and filters the artists.
	/// It takes the response writer, the request and the artists as parameters.
	/// It returns the sorted and filtered artists.

    research := r.URL.Query().Get("research") // Get the research parameter
    sortParam := r.URL.Query().Get("sort") // Get the sort parameter
	minYear := r.URL.Query().Get("minYear") // Get the minYear parameter
	maxYear := r.URL.Query().Get("maxYear") // Get the maxYear parameter



	var filteredArtists []Artist
	if minYear != "" && maxYear != "" { // Check if the minYear and maxYear parameters are not empty
		min, _ := strconv.Atoi(minYear)
		max, _ := strconv.Atoi(maxYear)
		for _, artist := range artists { // Loop through the artists
			if artist.CreationDate >= min && artist.CreationDate <= max { // Check if the artist creation date is between the min and max years
				filteredArtists = append(filteredArtists, artist)
			}
		}
		artists = filteredArtists
	}

    if research != "" { // Check if the research parameter is not empty
        for _, artist := range artists { // Loop through the artists
            if strings.Contains(strings.ToLower(artist.Name), strings.ToLower(research)) { // Check if the artist name contains the research
                filteredArtists = append(filteredArtists, artist)
            }
        }
        artists = filteredArtists
    }
	if sortParam != "" { // Check if the sort parameter is not empty 	
    	if sortParam == "A-Z" { // Check if the sort parameter is A-Z
        	sort.Slice(artists, func(i, j int) bool {
            	return artists[i].Name < artists[j].Name
        	})
    	} else if sortParam == "Z-A" { // Check if the sort parameter is Z-A
        	sort.Slice(artists, func(i, j int) bool {
            	return artists[i].Name > artists[j].Name
        	})
    	}
	}

    return artists
}

func SortDatesLocations(datesLocations map[string][]string) []LocationDates { 
	/// SortDatesLocations is a function that sorts the dates and locations.
	/// It takes the dates and locations as parameters.
	/// It returns the sorted dates and locations.

	var locationDatesSlice []LocationDates
	for location, dates := range datesLocations { // Loop through the dates and locations
		locationDatesSlice = append(locationDatesSlice, LocationDates{location, dates})
	}
	sort.Slice(locationDatesSlice, func(i, j int) bool { // Sort the dates and locations
		layout := "02-01-2006" // Set the layout
		date1, _ := time.Parse(layout, locationDatesSlice[i].Dates[0])
		date2, _ := time.Parse(layout, locationDatesSlice[j].Dates[0])
		return date1.Before(date2)
	})
	return locationDatesSlice
}

func GetCoordinates(location string) (string, string, error) {
	/// GetCoordinates is a function that gets the coordinates of the location.
	/// It takes the location as a parameter.
	/// It returns the latitude, longitude and an error.

	params := url.Values{} // Create a map to store the parameters
	location = strings.Replace(location, "-", ",", -1) // Replace the - with ,
	params.Set("q", location)
	params.Set("format", "json") // Set the format to json

	resp, err := client.Get("https://nominatim.openstreetmap.org/search?format=json&q=" + url.QueryEscape(location)) // Get the coordinates from the API
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
		return "", "", fmt.Errorf("no results found for location: %s", location)
	}

	return results[0].Latitude, results[0].Longitude, nil // Return the latitude and longitude
}

func CodeErreur(w http.ResponseWriter, r *http.Request, status int, message string) {
	/// CodeErreur is a function that returns an error message.
	/// It takes the response writer, the request, the status and the message as parameters.
	/// No return.

	colorRed := "\033[31m"

	if status == 404 { // Check if the status is 404
		http.Error(w, "404 not found", http.StatusNotFound)
		fmt.Println(string(colorRed), "[SERVER_ALERT] - 404 : File not found, or missing...", message)
		if status == 400 { // Check if the status is 400
			http.Error(w, "400 Bad request", http.StatusBadRequest)
			fmt.Println(string(colorRed), "[SERVER_ALERT] - 400 : Bad request", message)
		}
		if status == 500 { // Check if the status is 500
			http.Error(w, "500 Internal server error", http.StatusInternalServerError)
			fmt.Println(string(colorRed), "[SERVER_ALERT] - 500 : Internal server error", message)
		}

	}
}