package routes

import (
	"database/sql"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

//AuthorByID gets a specific author by his ID
func AuthorByID(db *sql.DB, c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return fiber.NewError(400, "Invalid Author ID")
	}
	row := db.QueryRow(`SELECT id, email, name, last_name, created_at, updated_at FROM authors WHERE id = $1`, id)
	authorQuery := author{}
	err = row.Scan(&authorQuery.ID, &authorQuery.Email, &authorQuery.Name, &authorQuery.LastName, &authorQuery.CreatedAt, &authorQuery.UpdatedAt)
	if err != nil {
		return fiber.NewError(400, "Author not found")
	}
	return c.JSON(authorQuery)
}
