package main

import (

	//"os"
	//"bufio"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"net/http"
)


var respart []string
var resploc []string
var respdat []string
var resprel []string
var err string


type Artist struct{
	id string
	image string
	name string
	members []string
	creationDate int
	firstAlbum string
	locations string
	concertDates string
	relations string
}
type Locations struct{
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
}/*
func convertion(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
		s[i] = string(b[i])
	}
	return string.Join(s,",")
}*/
func data_art() {
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
	
	var artist Artist
	err = json.Unmarshal(bodyart, &artist)
	fmt.Println(artist)
	
}

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
	
	
	var locations Locations
	err = json.Unmarshal(bodyloc, &locations)
	fmt.Println(locations)
	
	
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
	
	var dates Dates
	err = json.Unmarshal(bodydat, &dates)
	fmt.Println(dates)
	
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
	
	var relations Relations
	err = json.Unmarshal(bodyrel, &relations)
	fmt.Println(relations)
	
}

func main() {
	data_art()
	data_loc()
	data_dat()
	data_rel()
}