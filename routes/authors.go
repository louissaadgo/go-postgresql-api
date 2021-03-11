package routes

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type author struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	LastName  string `json:"lastName"`
	CreatedAt string `json:"createdAt"`
}

//Authors k
func Authors(db *sql.DB, c *fiber.Ctx) error {
	rows, err := db.Query(`SELECT * FROM authors;`)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Something went wrong.")
	}
	defer rows.Close()
	authors := make([]author, 0)
	for rows.Next() {
		auth := author{}
		err := rows.Scan(&auth.ID, &auth.Name, &auth.LastName, &auth.CreatedAt)
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
