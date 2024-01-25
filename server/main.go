package main

import (
	"fmt"
	"log"
	"net/http"
)

func HandlerMain(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "HTML/index.html")
}

func main() {

	http.HandleFunc("/", HandlerMain)

	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
