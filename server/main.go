package main

import (

	//"os"
	//"bufio"
	"io/ioutil"
	//"encoding/json"
	"fmt"
	"net/http"
	"log"
)

var response string
var respart []string
var resploc []string
var respdat []string
var resprel []string
var err []string
var name string

type Artist struct{
	ID string
	Image string
	Name string
	Members []string
	CreationDate int
	FirstAlbum string
}
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
	fmt.Println(string(bodyart))
	
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
	fmt.Println(string(bodyloc))
	
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
	fmt.Println(string(bodyrel))
	
}

func HandlerMain(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "HTML/index.html")
}

func main() {
	data_art()
	data_loc()
	data_dat()
	data_rel()

	http.HandleFunc("/", HandlerMain)

	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}