package main

import (
	"log"
	"net/http"
)

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
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

func (app *application) memberRegistrationFileHandler(w http.ResponseWriter, r *http.Request) {
	var reqBody MemberRegistrationFileRequest
	err := readJSON(r, &reqBody)
	if err != nil {
		log.Println(err)

		err = writeJSON(w, http.StatusBadRequest, envelope{"message": "invalid request"}, nil)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if app.reqFTPConn == nil {
		log.Println("reqFTPConn is nil")

		err = writeJSON(w, http.StatusInternalServerError, envelope{"message": "the server encountered a problem"}, nil)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	file, err := retrieveCSVFile(app.reqFTPConn, reqBody.File)
	if err != nil {
		log.Println("failed to retrieve file: ", err)

		err = writeJSON(w, http.StatusBadRequest, envelope{"message": "invalid request"}, nil)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	members, err := readCSVFile(file)
	if err != nil {
		log.Println("failed to read file: ", err)

		err = writeJSON(w, http.StatusBadRequest, envelope{"message": "invalid request"}, nil)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Register each member (include request data validation)
	// Collect the response of each member registration to be written to the response CSV file

	data := make([]MemberRegistrationResponse, 0)
	for _, member := range members {
		data = append(data, MemberRegistrationResponse{
			Name:   member.Name,
			Email:  member.Email,
			Phone:  member.Phone,
			Age:    member.Age,
			Active: member.Active,
		})
	}

	err = writeJSON(w, http.StatusCreated, envelope{"message": "members registered successfully", "data": data}, nil)
	if err != nil {
		log.Println(err)

		err = writeJSON(w, http.StatusInternalServerError, envelope{"message": "the server encountered a problem"}, nil)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Create & send response CSV file to FTP server
	// File name: <req-file-name>_response_timestamp.csv
	// Header: status (success/ error), response_body ({...})
}
