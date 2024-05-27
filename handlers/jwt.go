package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sangeet-server/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func GetToken(c *fiber.Ctx, user models.Token) (string, error) {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		log.Fatal("Please set JWT_KEY in you .env")
	}
	fmt.Println("key:", key)
	t := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{"user": user, "exp": time.Now().Add(72 * time.Hour).Unix()},
	)
	s, err := t.SignedString([]byte(key))
	if err != nil {
		return "", c.Status(http.StatusInternalServerError).
			JSON(fiber.Map{"error": "Something went wrong. Please try again.", "message": err.Error()})
	}
	return s, nil
}

// func VerifyToken(tokenString string) {
//     [
//     isOk := jwt.Parse()
// }

// func RefreshToken(tokenString) {
// }
