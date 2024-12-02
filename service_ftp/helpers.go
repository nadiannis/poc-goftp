package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/jlaffaye/ftp"
)

type envelope map[string]any

func writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	jsonBytes, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(jsonBytes)

	return nil
}

func readJSON(r *http.Request, dst any) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		return err
	}

	return nil
}

func connectToFTPServer(addr string, username string, password string) (*ftp.ServerConn, error) {
	conn, err := ftp.Dial(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to FTP server: %v", err)
	}

	err = conn.Login(username, password)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %v", err)
	}

	return conn, nil
}

func retrieveCSVFile(conn *ftp.ServerConn, filePath string) (file *ftp.Response, err error) {
	file, err = conn.Retr(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func readCSVFile(file *ftp.Response) (records []MemberRegistrationRequest, err error) {
	reader := csv.NewReader(file)
	defer file.Close()

	// Skip the header
	rowNum := 1
	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

	for {
		rowNum++
		row, err := reader.Read()
		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			return nil, err
		}

		age, err := strconv.Atoi(row[indexAge])
		if err != nil {
			return nil, fmt.Errorf("row %d: %v", rowNum, err)
		}

		active, err := strconv.ParseBool(row[indexActive])
		if err != nil {
			return nil, fmt.Errorf("row %d: %v", rowNum, err)
		}

		data := MemberRegistrationRequest{
			Name:   row[indexName],
			Email:  row[indexEmail],
			Phone:  row[indexPhone],
			Age:    age,
			Active: active,
		}
		fmt.Printf("%+v\n", data)

		records = append(records, data)
	}

	return records, nil
}

// func readCSVFile(file *ftp.Response) (contents [][]any, err error) {
// 	reader := csv.NewReader(file)
// 	defer file.Close()

// 	for {
// 		record, err := reader.Read()
// 		if err != nil {
// 			if err == io.EOF {
// 				err = nil
// 			}
// 			break
// 		}

// 		if len(dataHeaders) == 0 {
// 			dataHeaders = record
// 			continue
// 		}

// 		row := make([]any, 0)
// 		for _, value := range record {
// 			row = append(row, value)
// 		}

// 		contents = append(contents, row)
// 	}

// 	return contents, nil
// }

// func readCSVFile(file *ftp.Response) (contents [][]string, err error) {
// 	reader := csv.NewReader(file)
// 	defer file.Close()

// 	contents, err = reader.ReadAll()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return contents, nil
// }
