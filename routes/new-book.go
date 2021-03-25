package routes

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/louissaadgo/go-postgresql-api/auth"
)

//NewBook creates a new book
func NewBook(db *sql.DB, c *fiber.Ctx) error {
	bookNew := book{}
	err := c.BodyParser(&bookNew)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Invalid JSON")
	}
	if len(bookNew.Title) > 40 || len(bookNew.Title) == 0 {
		return fiber.NewError(400, "Book title must be 40 characters or less")
	}
	if len(bookNew.Description) < 40 || len(bookNew.Description) > 1000 {
		return fiber.NewError(400, "Book description must be 40 characters or more, 1000 at max")
	}
	row := db.QueryRow(`SELECT id, email FROM authors WHERE id = $1;`, bookNew.AuthorID)
	var authID int
	var email string
	row.Scan(&authID, &email)
	if authID == 0 {
		return fiber.NewError(400, "Author not found")
	}
	cookie := c.Cookies("session")
	maker, _ := auth.NewPasetoMaker(signkey)
	payload, err := maker.VerifyToken(cookie)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Invalid token")
	}
	newrow := db.QueryRow(`SELECT id FROM authors WHERE email = $1;`, payload.Username)
	var newID int
	newrow.Scan(&newID)
	if newID != bookNew.AuthorID {
		return fiber.NewError(400, "You must create a book under your ID")
	}
	_, err = db.Exec(`INSERT INTO books(title, author_id, description)
	VALUES($1, $2, $3);`, bookNew.Title, bookNew.AuthorID, bookNew.Description)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Something went wrong, please try again soon.")
	}
	return c.SendString("Book created successfully")
}
