package main

import (
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		err := writeJSON(w, http.StatusNotFound, envelope{"message": "the resource could not be found"}, nil)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	err := writeJSON(w, http.StatusOK, envelope{"message": "API is running on port 8080"}, nil)
	if err != nil {
		log.Println(err)

		err = writeJSON(w, http.StatusInternalServerError, envelope{"message": "the server encountered a problem"}, nil)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
