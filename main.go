package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var db *sql.DB
var err error

const port string = ":3024"

type author struct {
	Name     string `json:"name"`
	LastName string `json:"lastName"`
}

type book struct {
	Title    string `json:"title"`
	AuthorID int    `json:"authorID"`
}

func main() {
	db, err = sql.Open("postgres", "postgres://postgres:2400@localhost/library?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Connected to the database")
	defer db.Close()
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
	authorNew := author{}
	err := json.NewDecoder(r.Body).Decode(&authorNew)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	if len(authorNew.Name) > 26 || len(authorNew.LastName) > 26 {
		http.Error(w, "Author name and last name must be less than 26 characters", http.StatusBadRequest)
		return
	}
	_, err = db.Exec(`INSERT INTO authors(name, last_name)
	VALUES($1, $2);`, authorNew.Name, authorNew.LastName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong, please try again soon.", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprintln(w, "Author was successfully created")
}

func newBook(w http.ResponseWriter, r *http.Request) {

}
