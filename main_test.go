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
	mux.HandleFunc("/author/{id}", AuthorByID).Methods("GET")
	mux.HandleFunc("/book/{id}", BookByID).Methods("GET")
	return mux
}

func TestBookByID(t *testing.T) {
	DB, err = sql.Open("postgres", "postgres://postgres:2400@localhost/library?sslmode=disable")
	assert.Equal(t, nil, err, "Failed to connect to the database")
	defer DB.Close()
	req, _ := http.NewRequest("GET", "/book/1", nil)
	res := httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "Failed to handle CORRECT data")
	req, _ = http.NewRequest("GET", "/book/0", nil)
	res = httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 400, res.Code, "Failed to handle INCORRECT data")
	req, _ = http.NewRequest("GET", "/book/test", nil)
	res = httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 400, res.Code, "Failed to handle INCORRECT data")
	req, _ = http.NewRequest("GET", "/book/99999999", nil)
	res = httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 400, res.Code, "Failed to handle INCORRECT data")
}

func TestAuthorByID(t *testing.T) {
	DB, err = sql.Open("postgres", "postgres://postgres:2400@localhost/library?sslmode=disable")
	assert.Equal(t, nil, err, "Failed to connect to the database")
	defer DB.Close()
	req, _ := http.NewRequest("GET", "/author/1", nil)
	res := httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "Failed to handle CORRECT data")
	req, _ = http.NewRequest("GET", "/author/0", nil)
	res = httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 400, res.Code, "Failed to handle INCORRECT data")
	req, _ = http.NewRequest("GET", "/author/test", nil)
	res = httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 400, res.Code, "Failed to handle INCORRECT data")
	req, _ = http.NewRequest("GET", "/author/99999999", nil)
	res = httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 400, res.Code, "Failed to handle INCORRECT data")
}

func TestNewAuthor(t *testing.T) {
	DB, err = sql.Open("postgres", "postgres://postgres:2400@localhost/library?sslmode=disable")
	assert.Equal(t, nil, err, "Failed to connect to the database")
	defer DB.Close()
	payload := strings.NewReader(`{
    "name": "testName",
    "lastName": "testLastName"
}`)
	req, _ := http.NewRequest("POST", "/new/author", payload)
	res := httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "Failed to handle CORRECT data")
	payload = strings.NewReader(`{
    "name": "testName0000000000000000000000000",
    "lastName": "testLastName0000000000000000000000000"
}`)
	req, _ = http.NewRequest("POST", "/new/author", payload)
	res = httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 400, res.Code, "Failed to handle INCORRECT data")
}

func TestNewBook(t *testing.T) {
	DB, err = sql.Open("postgres", "postgres://postgres:2400@localhost/library?sslmode=disable")
	assert.Equal(t, nil, err, "Failed to connect to the database")
	defer DB.Close()
	payload := strings.NewReader(`{
		"title": "testTitle",
		"authorID": 1
	}`)
	req, _ := http.NewRequest("POST", "/new/book", payload)
	res := httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "Failed to handle CORRECT data")
	payload = strings.NewReader(`{
    "title": "testTitle000000000000000000000000000000000000000000",
    "authorID": 999999
}`)
	req, _ = http.NewRequest("POST", "/new/book", payload)
	res = httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 400, res.Code, "Failed to handle INCORRECT data")
}
