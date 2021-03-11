package routes

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

//BookByID j
func BookByID(db *sql.DB, c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return fiber.NewError(400, "Invalid Book ID")
	}
	row := db.QueryRow(`SELECT * FROM books WHERE id = $1`, id)
	bookQuery := book{}
	err = row.Scan(&bookQuery.ID, &bookQuery.Title, &bookQuery.AuthorID, &bookQuery.PublishedAt)
	if err != nil {
		return fiber.NewError(400, "Book not found")
	}
	return c.JSON(bookQuery)
}
