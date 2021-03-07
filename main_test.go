package main

import (
	"github.com/gorilla/mux"
)

func router() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/new/author", NewAuthor).Methods("POST")
	mux.HandleFunc("/new/book", NewBook).Methods("POST")
	return mux
}
