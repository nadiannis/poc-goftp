package main

const (
	indexName = iota
	indexEmail
	indexPhone
	indexAge
	indexActive
)

type MemberRegistrationFileRequest struct {
	File string `json:"file"`
}

type MemberRegistrationRequest struct {
	Name   string
	Email  string
	Phone  string
	Age    int
	Active bool
}

type MemberRegistrationResponse struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Age    int    `json:"age"`
	Active bool   `json:"active"`
}
