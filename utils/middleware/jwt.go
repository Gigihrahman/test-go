package middleware

import (
	"log"
	"net/http"
	"os"

	"test-rakamin/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.ErrorResponseFiber(c, http.StatusUnauthorized, "Unauthorized", "Token is required")
		}

		tokenString := authHeader

		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			log.Printf("JWT parsing error: %v", err)
			return utils.ErrorResponseFiber(c, http.StatusUnauthorized, "Unauthorized", "Invalid or expired token")
		}

		userIDFloat, ok := claims["id"].(float64)
		if !ok {
			return utils.ErrorResponseFiber(c, http.StatusUnauthorized, "Invalid token claims", "User ID not found in token")
		}
		c.Locals("user_id", uint(userIDFloat))

		return c.Next()
	}
}
