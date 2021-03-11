package routes

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

//NewAuthor n
func NewAuthor(db *sql.DB, c *fiber.Ctx) error {
	authorNew := author{}
	err := c.BodyParser(&authorNew)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Invalid JSON.")
	}
	if len(authorNew.Name) > 25 || len(authorNew.LastName) > 25 {
		return fiber.NewError(400, "Author name and last name must be 25 characters or less")
	}
	if len(authorNew.Name) == 0 || len(authorNew.LastName) == 0 {
		return fiber.NewError(400, "Author name and last name must be at least 1 character")
	}
	_, err = db.Exec(`INSERT INTO authors(name, last_name)
	VALUES($1, $2);`, authorNew.Name, authorNew.LastName)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Something went wrong, please try again soon.")
	}
	return c.SendString("Author created successfully")
}
