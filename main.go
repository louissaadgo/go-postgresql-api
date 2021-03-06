package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB

const port string = ":3024"

func init() {
	db, err := sql.Open("postgres", "postgres://postgres:2400@localhost/library?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Connected to the database")
	defer db.Close()
}

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/books", books).Methods("GET")
	mux.HandleFunc("/authors", authors).Methods("GET")
	mux.HandleFunc("/new/author", newAuthor).Methods("POST")
	mux.HandleFunc("/new/book", newBook).Methods("POST")
	log.Fatal(http.ListenAndServe(port, mux))
}

func books(w http.ResponseWriter, r *http.Request) {

}

func authors(w http.ResponseWriter, r *http.Request) {

}

func newAuthor(w http.ResponseWriter, r *http.Request) {

}

func newBook(w http.ResponseWriter, r *http.Request) {

}
