package routes

import (
	"database/sql"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type queryBooks struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
	Name        string `json:"name"`
	LastName    string `json:"lastName"`
}

//Books shows all the books
func Books(db *sql.DB, c *fiber.Ctx) error {
	rows, err := db.Query(`SELECT books.id, title, published_at, authors.name, authors.last_name FROM books
	JOIN authors ON authors.id = books.author_id;`)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Something went wrong.")
	}
	defer rows.Close()
	books := make([]queryBooks, 0)
	for rows.Next() {
		bookNew := queryBooks{}
		err := rows.Scan(&bookNew.ID, &bookNew.Title, &bookNew.PublishedAt, &bookNew.Name, &bookNew.LastName)
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
	c.JSON(books)
	return nil
}
