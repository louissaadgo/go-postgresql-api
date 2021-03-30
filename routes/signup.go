package routes

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
		fmt.Println(err)
		return fiber.NewError(400, "Invalid JSON.")
	}
	errors := authorNew.Validate()
	if errors != nil {
		c.Status(400)
		c.JSON(errors)
		return nil
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
