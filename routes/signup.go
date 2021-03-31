package routes

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/louissaadgo/go-postgresql-api/auth"
	"github.com/louissaadgo/go-postgresql-api/validation"
)

//Signup creates a new author
func Signup(db *sql.DB, c *fiber.Ctx) error {
	authorNew := validation.NewAuthor{}
	err := c.BodyParser(&authorNew)
	if err != nil {
		return fiber.NewError(400, "Invalid JSON.")
	}
	errors := authorNew.Validate()
	if errors != nil {
		c.Status(400)
		c.JSON(errors)
		return nil
	}
	_, err = db.Exec(`INSERT INTO authors(email, password, name, last_name)
	VALUES($1, $2, $3, $4);`, authorNew.Email, auth.HashPassword(authorNew.Password), authorNew.Name, authorNew.LastName)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Something went wrong, please try again soon.")
	}
	cookie, err := auth.NewPasetoCookie(signkey, authorNew.Email, time.Hour*24)
	if err != nil {
		return fiber.NewError(400, "Something went wrong, please try again soon.")
	}
	c.Cookie(cookie)
	return c.SendString("Author created successfully")
}
