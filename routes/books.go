package routes

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

//Books gets all the books in the database
func Books(db *sql.DB, c *fiber.Ctx) error {
	rows, err := db.Query(`SELECT books.id, title, description, published_at, authors.name, authors.last_name FROM books
	JOIN authors ON authors.id = books.author_id;`)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Something went wrong, please try again soon.")
	}
	defer rows.Close()
	books := make([]queryBooks, 0)
	for rows.Next() {
		bookNew := queryBooks{}
		err := rows.Scan(&bookNew.ID, &bookNew.Title, &bookNew.Description, &bookNew.PublishedAt, &bookNew.Name, &bookNew.LastName)
		if err != nil {
			fmt.Println(err)
			return fiber.NewError(400, "Something went wrong.")
		}
		books = append(books, bookNew)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Something went wrong.")
	}
	return c.JSON(books)
}
