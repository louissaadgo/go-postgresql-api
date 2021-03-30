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

//Signin signs an old author
func Signin(db *sql.DB, c *fiber.Ctx) error {
	authorOld := validation.OldAuthor{}
	err := c.BodyParser(&authorOld)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Invalid JSON.")
	}
	errors := authorOld.Validate()
	if errors != nil {
		c.Status(400)
		c.JSON(errors)
		return nil
	}
	h := sha256.Sum256([]byte(authorOld.Password))
	pass := hex.EncodeToString(h[:])
	row := db.QueryRow(`SELECT id, email, name, last_name, password FROM authors WHERE email = $1;`, authorOld.Email)
	authorQuery := author{}
	err = row.Scan(&authorQuery.ID, &authorQuery.Email, &authorQuery.Name, &authorQuery.LastName, &authorQuery.Password)
	if err != nil {
		return fiber.NewError(400, "User not found")
	}
	if pass != authorQuery.Password {
		return fiber.NewError(400, "Wrong credentials")
	}
	maker, _ := auth.NewPasetoMaker(signkey)
	token, _ := maker.CreateToken(authorQuery.Email, time.Hour*24)
	cookie := new(fiber.Cookie)
	cookie.Name = "session"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	c.Cookie(cookie)
	return c.SendString("Logged in successfully")
}
