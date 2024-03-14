package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mod/handlers"
)

// Middleware untuk memeriksa token JWT
func Authenticated(c *fiber.Ctx) error {
	// Ambil token dari header Authorization
	authHeader := c.Get("Authorization")

	// Periksa jika Authorization header tidak ada
	if authHeader == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Authorization header missing",
		})
	}

	// Pecah string menjadi bagian "Bearer" dan token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid Authorization header format",
		})
	}
	tokenString := parts[1]

	// Validasi token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi metode penandatanganan token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return handlers.SecretKey, nil
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Unauthorized",
		})
	}

	// Periksa apakah token valid
	if !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token",
		})
	}

	// Lanjutkan jika token valid dan belum kedaluwarsa
	return c.Next()
}
