package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/scan/", func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		log.Print(url)
		fmt.Fprintf(w, "This is a response %v", url)
	})

	log.Print("Server is starting")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
