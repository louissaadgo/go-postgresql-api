package routes

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

//Authors gets all the authors in the database
func Authors(db *sql.DB, c *fiber.Ctx) error {
	rows, err := db.Query(`SELECT id, email, name, last_name, created_at, updated_at FROM authors;`)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Something went wrong, please try again soon.")
	}
	defer rows.Close()
	authors := make([]author, 0)
	for rows.Next() {
		auth := author{}
		err := rows.Scan(&auth.ID, &auth.Email, &auth.Name, &auth.LastName, &auth.CreatedAt, &auth.UpdatedAt)
		if err != nil {
			fmt.Println(err)
			return fiber.NewError(400, "Something went wrong.")
		}
		authors = append(authors, auth)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Something went wrong.")
	}
	return c.JSON(authors)
}
