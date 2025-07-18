package middleware

import (
	"backend/database"
	"backend/prisma/db"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(jwtSecret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := getTokenFromRequest(c)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization token missing"})
		}

		claims, err := validateToken(tokenString, jwtSecret)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		sessionUUID, ok := claims["sessionUUID"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid session UUID in token"})
		}

		ctx := context.Background()
		session, err := database.Client().ActiveSessions.FindUnique(
			db.ActiveSessions.SessionUUID.Equals(sessionUUID),
		).Exec(ctx)
		if err == db.ErrNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Session not found"})
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
		}

		if time.Now().After(session.Exp) {
			database.Client().ActiveSessions.FindUnique(
				db.ActiveSessions.SessionUUID.Equals(sessionUUID),
			).Delete().Exec(ctx)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Session expired"})
		}

		c.Locals(("userID"), claims["userID"])
		c.Locals("sessionUUID", sessionUUID)
		return c.Next()
	}
}

func getTokenFromRequest(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}
	return c.Cookies("auth")
}

func validateToken(tokenString string, jwtSecret []byte) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
