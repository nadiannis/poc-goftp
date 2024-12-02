package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jlaffaye/ftp"
	"github.com/joho/godotenv"
)

type application struct {
	reqFTPConn *ftp.ServerConn
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	reqFTPConn, err := connectToFTPServer(os.Getenv("REQ_FTP_ADDR"), os.Getenv("REQ_FTP_USERNAME"), os.Getenv("REQ_FTP_PASSWORD"))
	if err != nil {
		log.Println(err)
	}

	app := application{
		reqFTPConn: reqFTPConn,
	}

	http.HandleFunc("GET /", app.homeHandler)
	http.HandleFunc("POST /members/registration/files", app.memberRegistrationFileHandler)

	log.Println("API server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
