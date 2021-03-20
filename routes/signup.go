package routes

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
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
	claims := myCustomClaims{
		name: authorNew.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: 15000,
			Issuer:    authorNew.Name,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, _ := token.SignedString(signkey)
	return c.SendString("Author created successfully" + ss)
}
