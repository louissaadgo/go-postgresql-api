package main

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func router() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/new/author", NewAuthor).Methods("POST")
	mux.HandleFunc("/new/book", NewBook).Methods("POST")
	return mux
}

func TestNewAuthor(t *testing.T) {
	DB, err = sql.Open("postgres", "postgres://postgres:2400@localhost/library?sslmode=disable")
	assert.Equal(t, nil, err, "Failed to connect to the database")
	defer DB.Close()
	payload := strings.NewReader(`{` + "" + `
    "name": "testName",` + "" + `
    "lastName": "testLastName"` + "" + `
}`)
	req, _ := http.NewRequest("POST", "/new/author", payload)
	res := httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "Failed to handle CORRECT data")
	payload = strings.NewReader(`{` + "" + `
    "name": "TeeeeeeessssssssssssttttttttttNAMEeeeeee",` + "" + `
    "lastName": "TestlassssssssssssssssssssstNaaaaaaaaaaaaaaame"` + "" + `
}`)
	req, _ = http.NewRequest("POST", "/new/author", payload)
	res = httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 400, res.Code, "Failed to handle INCORRECT data")
}
