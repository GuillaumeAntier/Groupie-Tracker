package main

import (
	"fmt"
	"net/http"
	"GroupieTracker"
	"log"
)


func main() {
	fmt.Println(string("\033[34m"), "[SERVER_INFO] : Starting local Server...")
	
	http.Handle("/JS/", http.StripPrefix("/JS/", http.FileServer(http.Dir("JS"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/presentation", GroupieTracker.ServeHome)
	http.HandleFunc("/index", GroupieTracker.HandlerMain)
	http.HandleFunc("/result", GroupieTracker.SearchHandler)
	http.HandleFunc("/event", GroupieTracker.EventHandler)

	fmt.Println(string("\033[32m"), "[SERVER_READY] : on http://localhost:8080/index ✅ ")
	fmt.Println(string("\033[33m"), "[SERVER_INFO] : To stop the program : Ctrl + c")
	log.Fatal(http.ListenAndServe(":8080", nil))

}