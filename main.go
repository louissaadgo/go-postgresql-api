package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/louissaadgo/go-postgresql-api/routes"
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
	app := fiber.New()
	app.Get("/books", func(c *fiber.Ctx) error {
		return routes.Books(DB, c)
	})
	app.Get("/authors", func(c *fiber.Ctx) error {
		return routes.Authors(DB, c)
	})
	app.Listen(port)
	// mux.HandleFunc("/author/{id}", AuthorByID).Methods("GET")
	// mux.HandleFunc("/book/{id}", BookByID).Methods("GET")
	// mux.HandleFunc("/new/author", NewAuthor).Methods("POST")
	// mux.HandleFunc("/new/book", NewBook).Methods("POST")
	// log.Fatal(http.ListenAndServe(port, mux))
}

//BookByID shows a specific book
func BookByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil || id < 1 {
		http.Error(w, "Invalid Book ID", http.StatusBadRequest)
		return
	}
	row := DB.QueryRow(`SELECT * FROM books WHERE id = $1`, id)
	bookQuery := book{}
	err = row.Scan(&bookQuery.ID, &bookQuery.Title, &bookQuery.AuthorID, &bookQuery.PublishedAt)
	if err != nil {
		http.Error(w, "Book not found", http.StatusBadRequest)
		return
	}
	bs, _ := json.Marshal(bookQuery)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	fmt.Fprintln(w, string(bs))
}

//AuthorByID shows a specific author
func AuthorByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil || id < 1 {
		http.Error(w, "Invalid Author ID", http.StatusBadRequest)
		return
	}
	row := DB.QueryRow(`SELECT * FROM authors WHERE id = $1`, id)
	authorQuery := author{}
	err = row.Scan(&authorQuery.ID, &authorQuery.Name, &authorQuery.LastName, &authorQuery.CreatedAt)
	if err != nil {
		http.Error(w, "Author not found", http.StatusBadRequest)
		return
	}
	bs, _ := json.Marshal(authorQuery)
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
		http.Error(w, "Something went wrong, please try again soon.", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, "Book was successfully created")
}
