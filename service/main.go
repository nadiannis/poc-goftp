package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("GET /", homeHandler)

	log.Println("API server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
