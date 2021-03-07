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

//DB connection
var DB *sql.DB
var err error

const port string = ":3024"

type author struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	LastName  string `json:"lastName"`
	CreatedAt string `json:"createdAt"`
}

type book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	AuthorID    int    `json:"authorID"`
	PublishedAt string `json:"publishedAt"`
}

type queryBooks struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
	Name        string `json:"name"`
	LastName    string `json:"lastName"`
}

func main() {
	DB, err = sql.Open("postgres", "postgres://postgres:2400@localhost/library?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Connected to the database")
	defer DB.Close()
	mux := mux.NewRouter()
	mux.HandleFunc("/books", Books).Methods("GET")
	mux.HandleFunc("/authors", Authors).Methods("GET")
	mux.HandleFunc("/new/author", NewAuthor).Methods("POST")
	mux.HandleFunc("/new/book", NewBook).Methods("POST")
	log.Fatal(http.ListenAndServe(port, mux))
}

//Books shows all the books
func Books(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`SELECT books.id, title, published_at, authors.name, authors.last_name FROM books
	JOIN authors ON authors.id = books.author_id;`)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusBadRequest)
		return
	}
	defer rows.Close()
	books := make([]queryBooks, 0)
	for rows.Next() {
		bookNew := queryBooks{}
		err := rows.Scan(&bookNew.ID, &bookNew.Title, &bookNew.PublishedAt, &bookNew.Name, &bookNew.LastName)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Something went wrong.", http.StatusBadRequest)
			return
		}
		books = append(books, bookNew)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusBadRequest)
		return
	}
	bs, _ := json.Marshal(books)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	fmt.Fprintln(w, string(bs))
}

//Authors shows all the authors
func Authors(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query(`SELECT * FROM authors;`)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusBadRequest)
		return
	}
	defer rows.Close()
	authors := make([]author, 0)
	for rows.Next() {
		auth := author{}
		err := rows.Scan(&auth.ID, &auth.Name, &auth.LastName, &auth.CreatedAt)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Something went wrong.", http.StatusBadRequest)
			return
		}
		authors = append(authors, auth)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong.", http.StatusBadRequest)
		return
	}
	bs, _ := json.Marshal(authors)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	fmt.Fprintln(w, string(bs))
}

//NewAuthor adds a new author
func NewAuthor(w http.ResponseWriter, r *http.Request) {
	authorNew := author{}
	err := json.NewDecoder(r.Body).Decode(&authorNew)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if len(authorNew.Name) > 25 || len(authorNew.LastName) > 25 {
		http.Error(w, "Author name and last name must be 25 characters or less", http.StatusBadRequest)
		return
	}
	if len(authorNew.Name) == 0 || len(authorNew.LastName) == 0 {
		http.Error(w, "Author name and last name must be at least 1 character", http.StatusBadRequest)
		return
	}
	_, err = DB.Exec(`INSERT INTO authors(name, last_name)
	VALUES($1, $2);`, authorNew.Name, authorNew.LastName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong, please try again soon.", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "Author was successfully created")
}

//NewBook adds a new book
func NewBook(w http.ResponseWriter, r *http.Request) {
	bookNew := book{}
	err := json.NewDecoder(r.Body).Decode(&bookNew)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if len(bookNew.Title) > 40 || len(bookNew.Title) == 0 {
		http.Error(w, "Book title must be 40 characters or less", http.StatusBadRequest)
		return
	}
	row := DB.QueryRow(`SELECT id FROM authors WHERE id = $1;`, bookNew.AuthorID)
	var authID int
	row.Scan(&authID)
	if authID == 0 {
		http.Error(w, "Author not found", http.StatusBadRequest)
		return
	}
	_, err = DB.Exec(`INSERT INTO books(title, author_id)
	VALUES($1, $2);`, bookNew.Title, bookNew.AuthorID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong, please try again soon.", http.StatusServiceUnavailable)
		return
	}
	fmt.Fprintln(w, "Book was successfully created")
}
