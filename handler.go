package GroupieTracker

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
)

func ServePageArtist(w http.ResponseWriter, r *http.Request, html string, data []Artist) {
	///ServePageArtist is a function that serves the page artist.html with the data of the artists. 
	///It takes the response writer, the request and the html file as parameters. 
	///It returns the page with the data of the artists. 
	if r.URL.Path == "/index" {
	} else {
		CodeErreur(w, r, 404, "Page not found")
		return
	}

	data = sortAlphabetically(w,r,data) //Sort the artists alphabetically

	page, err := template.ParseFiles("HTML/" + html)
	if err != nil {
		CodeErreur(w, r, 500, "Template not found : index.html")
		return
	}
	err = page.Execute(w, data)
	if err != nil {
		CodeErreur(w, r, 500, "Internal server error")
		return
	}
}
func ServePagePresentation(w http.ResponseWriter, r *http.Request, html string) {
	///ServePagePresentation is a function that serves the page presentation.html.
	///It takes the response writer, the request and the html file as parameters.
	///It returns the page presentation.html.
	http.ServeFile(w, r, "HTML/"+html)
	}

func ServePageResult(w http.ResponseWriter, r *http.Request, html string, data []Artist) {
	///ServePageResult is a function that serves the page result.html with the data of the artists.
	///It takes the response writer, the request and the html file as parameters.
	///It returns the page with the data of the artists.
	if r.URL.Path != "/result" {
		CodeErreur(w, r, 404, "Page not found")
		return
	}

	data = SortAndFilter(w,r,data)

	page, err := template.ParseFiles("HTML/" + html)
	if err != nil {
		CodeErreur(w, r, 500, "Template not found : result.html")
		return
	}
	err = page.Execute(w, data)
	if err != nil {
		CodeErreur(w, r, 500, "Internal server error")
		return
	}
}

func ServePageEvent(w http.ResponseWriter, r *http.Request, html string, dataArtist []Artist ,data []LocationDates, coordinatesMap map[string][]LocationCity) {
	///ServePageEvent is a function that serves the page event.html with the data of the artists, the dates and the locations.
	///It takes the response writer, the request and the html file as parameters.
	///It returns the page with the data of the artists, the dates and the locations.
	if r.URL.Path != "/event" {
		CodeErreur(w, r, 404, "Page not found")
		return
	}
	page, err := template.ParseFiles("HTML/" + html)
	if err != nil {
		CodeErreur(w, r, 500, "Template not found : event.html")
		return
	}
	// Execute the template
	err = page.Execute(w, struct { // Create a struct to pass the data to the template
		Artists        []Artist
		DatesLocations []LocationDates
		Coordinates    map[string][]LocationCity
	}{	// Pass the data to the template
		Artists:        dataArtist,
		DatesLocations: data,
		Coordinates:    coordinatesMap,
	})
	if err != nil {
		CodeErreur(w, r, 500, "Internal server error")
		return
	}
}



func ServeHome(w http.ResponseWriter, r *http.Request) {
	///ServeHome is a function that serves the page presentation.html.
	///It takes the response writer and the request as parameters.
	///It returns the page presentation.html.
	ServePagePresentation(w, r, "presentation.html")

}



func HandlerMain(w http.ResponseWriter, r *http.Request) {
	///HandlerMain is a function that serves the page index.html with the data of the artists.
	///It takes the response writer and the request as parameters.
	///It returns the page with the data of the artists.
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyart, err := OpenAPI(w, r, url+"artists")
	if err != nil {
		CodeErreur(w, r, 500, "Internal server error")
		return
	}
	var artist []Artist
	err = json.Unmarshal(bodyart, &artist)
	if err != nil {
		CodeErreur(w, r, 500, "Internal server error")
		return
	}
	ServePageArtist(w, r, "index.html", artist)
}

func SearchHandler(w http.ResponseWriter, r *http.Request) {
	///SearchHandler is a function that serves the page result.html with the data of the artists.
	///It takes the response writer and the request as parameters.
	///It returns the page with the data of the artists.
	research := r.URL.Query().Get("research")
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyart, err := OpenAPI(w, r, url+"artists")
	if err != nil {
		CodeErreur(w, r, 500, "Internal server error")
		return
	}
	var artist []Artist
	err = json.Unmarshal(bodyart, &artist)
	if err != nil {
		CodeErreur(w, r, 500, "Internal server error")
		return
	}
	// Filter the artists by the letter
	artist = FilterArtistsByLetter(artist, research)
	ServePageResult(w, r, "result.html", artist)
}

func EventHandler(w http.ResponseWriter, r *http.Request) {
	///EventHandler is a function that serves the page event.html with the data of the artists, the dates and the locations.
	///It takes the response writer and the request as parameters.
	///It returns the page with the data of the artists, the dates and the locations.
	id := r.URL.Query().Get("id")
	url := "https://groupietrackers.herokuapp.com/api/"
	bodyevent, err := OpenAPI(w, r, url+"relation/"+id)
	if err != nil {
		CodeErreur(w, r, 500, "Internal server error")
		return
	}
	var event Event
	err = json.Unmarshal(bodyevent, &event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	bodyart, err := OpenAPI(w, r, url+"artists")
	if err != nil {
		CodeErreur(w, r, 500, "Internal server error")
		return
	}
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
			CodeErreur(w, r, 500, "Internal server error")
			return
		}
		if a.ID == id {
			artist = append(artist, a)
		}
	}
	// Sort the dates and locations
	locationDateSlice := SortDatesLocations(event.DatesLocations)
	cordinatesMap := make(map[string][]LocationCity) // Create a map to store the coordinates

	for location := range event.DatesLocations { // Get the coordinates of the locations
		latitude, longitude, err := GetCoordinates(location)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cordinatesMap[location] = []LocationCity{{latitude, longitude}} // Store the coordinates in the map
	}
	ServePageEvent(w, r, "event.html", artist ,locationDateSlice, cordinatesMap)
}