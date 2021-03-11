package routes

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

//AuthorByID a
func AuthorByID(db *sql.DB, c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return fiber.NewError(400, "Invalid Author ID")
	}
	row := db.QueryRow(`SELECT * FROM authors WHERE id = $1`, id)
	authorQuery := author{}
	err = row.Scan(&authorQuery.ID, &authorQuery.Name, &authorQuery.LastName, &authorQuery.CreatedAt)
	if err != nil {
		return fiber.NewError(400, "Author not found")
	}
	return c.JSON(authorQuery)
}
