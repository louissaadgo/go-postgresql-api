package routes

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/louissaadgo/go-postgresql-api/auth"
	"github.com/louissaadgo/go-postgresql-api/validation"
)

//Signin signs an old author
func Signin(db *sql.DB, c *fiber.Ctx) error {
	authorOld := validation.OldAuthor{}
	err := c.BodyParser(&authorOld)
	if err != nil {
		return fiber.NewError(400, "Invalid JSON.")
	}
	errors := authorOld.Validate()
	if errors != nil {
		c.Status(400)
		c.JSON(errors)
		return nil
	}
	row := db.QueryRow(`SELECT id, email, name, last_name, password FROM authors WHERE email = $1;`, authorOld.Email)
	authorQuery := author{}
	err = row.Scan(&authorQuery.ID, &authorQuery.Email, &authorQuery.Name, &authorQuery.LastName, &authorQuery.Password)
	if err != nil {
		return fiber.NewError(400, "User not found")
	}
	if auth.HashPassword(authorOld.Password) != authorQuery.Password {
		return fiber.NewError(400, "Wrong credentials")
	}
	cookie, err := auth.NewPasetoCookie(signkey, authorQuery.Email, time.Hour*24)
	if err != nil {
		return fiber.NewError(400, "Something went wrong, please try again soon.")
	}
	c.Cookie(cookie)
	return c.SendString("Logged in successfully")
}
