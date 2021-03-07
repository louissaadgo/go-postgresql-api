package main

import (
	"database/sql"
	"fmt"
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
	assert.Equal(t, nil, err, "OK expected")
	fmt.Println("Connected to the database")
	defer DB.Close()
	payload := strings.NewReader(`{` + "" + `
    "name": "kfk",` + "" + `
    "lastName": "khalil"` + "" + `
}`)
	req, _ := http.NewRequest("POST", "/new/author", payload)
	res := httptest.NewRecorder()
	router().ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "OK expected")
}
