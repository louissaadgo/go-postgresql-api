package routes

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/louissaadgo/go-postgresql-api/auth"
)

//Signup creates a new author
func Signup(db *sql.DB, c *fiber.Ctx) error {
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
	if len(authorNew.Email) < 10 || len(authorNew.Email) > 200 {
		return fiber.NewError(400, "Author email must be at least 10 characters or 200 at max")
	}
	if len(authorNew.Password) < 8 || len(authorNew.Password) > 256 {
		return fiber.NewError(400, "Author password must be at least 8 characters or 256 at max")
	}
	h := sha256.Sum256([]byte(authorNew.Password))
	_, err = db.Exec(`INSERT INTO authors(email, password, name, last_name)
	VALUES($1, $2, $3, $4);`, authorNew.Email, hex.EncodeToString(h[:]), authorNew.Name, authorNew.LastName)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Something went wrong, please try again soon.")
	}
	maker, _ := auth.NewPasetoMaker(signkey)
	token, _ := maker.CreateToken(authorNew.Email, time.Hour*24)
	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.Cookie(cookie)
	return c.SendString("Author created successfully")
}
