package main

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/louissaadgo/go-postgresql-api/routes"
)

const port string = ":3024"

var (
	//DB connection
	DB  *sql.DB
	err error
)

func main() {
	DB, err = sql.Open("postgres", "postgres://postgres:2400@localhost/library?sslmode=disable")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer DB.Close()

	app := fiber.New()

	app.Get("/books", func(c *fiber.Ctx) error {
		return routes.Books(DB, c)
	})

	app.Get("/authors", func(c *fiber.Ctx) error {
		return routes.Authors(DB, c)
	})

	app.Post("/Signup", func(c *fiber.Ctx) error {
		return routes.Signup(DB, c)
	})

	app.Post("/new/book", func(c *fiber.Ctx) error {
		return routes.NewBook(DB, c)
	})

	app.Get("/author/:id", func(c *fiber.Ctx) error {
		return routes.AuthorByID(DB, c)
	})

	app.Get("/book/:id", func(c *fiber.Ctx) error {
		return routes.BookByID(DB, c)
	})

	app.Listen(port)
}
