package main

import (

	//"os"
	//"bufio"
	"html/template"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)
const port = ":8080"

var response string
var respart []string
var resploc []string
var respdat []string
var resprel []string
var err []string
var name string
var artist Artist

type Artist struct{
	ID string
	Image string
	//Name string
	//Members []string
	//CreationDate int
	//FirstAlbum string
}
/*type Locations struct{
	index []string
	id string
	locations string
	dates string
}
type Dates struct{
	index []string
	id string
	dates string
}
type Relations struct{
	index []string
	id string
	datesLocations string
}*//*
func data_art(w http.ResponseWriter, r *http.Request) {
	urlart := "https://groupietrackers.herokuapp.com/api/artists"
	respart, err := http.Get(urlart)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return
	}
	defer respart.Body.Close()
	bodyart, err := ioutil.ReadAll(respart.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture:", err)
		return 
	}
	//fmt.Println(string(bodyart))
	var artist []Artist
	err = json.Unmarshal(bodyart, &artist)
	fmt.Println(artist)
	servePage(w, r, "index.html", artist)
	
}*/

func data_loc() {
	urlloc := "https://groupietrackers.herokuapp.com/api/locations"
	resploc, err := http.Get(urlloc)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return
	}
	defer resploc.Body.Close()
	bodyloc, err := ioutil.ReadAll(resploc.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture:", err)
		return 
	}
	
	fmt.Println(string(bodyloc))
	//var locations []Locations
	//err = json.Unmarshal(bodyloc, &locations)
	//fmt.Println(string(locations))
	
	
}
func data_dat() {
	urldat := "https://groupietrackers.herokuapp.com/api/dates"
	respdat, err := http.Get(urldat)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return
	}
	defer respdat.Body.Close()
	bodydat, err := ioutil.ReadAll(respdat.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture:", err)
		return 
	}
	
	//var dates []Dates
	//err = json.Unmarshal(bodydat, &dates)
	//fmt.Println(string(dates))
	fmt.Println(string(bodydat))
	
}
func data_rel() {
	urlrel := "https://groupietrackers.herokuapp.com/api/relation"
	resprel, err := http.Get(urlrel)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return
	}
	defer resprel.Body.Close()
	bodyrel, err := ioutil.ReadAll(resprel.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture:", err)
		return 
	}
	//var relations []Relations
	//err = json.Unmarshal(bodyrel, &relations)
	//fmt.Println(string(relations))
	fmt.Println(string(bodyrel))
	
}

func servePage (w http.ResponseWriter, r *http.Request, html string, data []Artist) {
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
	urlart := "https://groupietrackers.herokuapp.com/api/artists"
	respart, err := http.Get(urlart)
	if err != nil {
		fmt.Println("Erreur lors de la requête HTTP:", err)
		return
	}
	defer respart.Body.Close()
	bodyart, err := ioutil.ReadAll(respart.Body)
	if err != nil {
		fmt.Println("Erreur lors de la lecture:", err)
		return 
	}
	var artist []Artist
	err = json.Unmarshal(bodyart, &artist)
	if err != nil {
		fmt.Println("Erreur lors de la lecture:", err)
		return
	}
	servePage(w, r, "index.html", artist)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	research := r.URL.Query().Get("research")
	fmt.Println(research)

	
	
	
}

func main() {
	//data_art()
	data_loc()
	data_dat()
	data_rel()

	http.HandleFunc("/", HandlerMain)
	http.HandleFunc("/result", getHandler)

	fmt.Println("Server is listening...")
	fmt.Println("(http://localhost:8080) - Serveur démarré sur le port", port)
	log.Fatal(http.ListenAndServe(":8080", nil))

}