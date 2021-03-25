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

//Signin signs an old author
func Signin(db *sql.DB, c *fiber.Ctx) error {
	authorNew := author{}
	err := c.BodyParser(&authorNew)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(400, "Invalid JSON.")
	}
	if len(authorNew.Email) < 10 || len(authorNew.Email) > 200 {
		return fiber.NewError(400, "Author email must be at least 10 characters or 200 at max")
	}
	if len(authorNew.Password) < 8 || len(authorNew.Password) > 256 {
		return fiber.NewError(400, "Author password must be at least 8 characters or 256 at max")
	}
	h := sha256.Sum256([]byte(authorNew.Password))
	pass := hex.EncodeToString(h[:])
	row := db.QueryRow(`SELECT id, email, name, last_name, password FROM authors WHERE email = $1;`, authorNew.Email)
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
