package routes

import (
	"backend/prisma/db"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx, client *db.PrismaClient, jwtSecret []byte) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Check if user exists
	ctx := context.Background()
	existingUser, err := client.User.FindUnique(
		db.User.Email.Equals(input.Email),
	).Exec(ctx)
	if err == nil && existingUser != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
	}
	if err != nil && err != db.ErrNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot hash password"})
	}

	// Create user
	user, err := client.User.CreateOne(
		db.User.Email.Set(input.Email),
		db.User.Password.Set(string(hashedPassword)),
	).Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot create user"})
	}

	sessionUUID := uuid.New()

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.ID
	claims["sessionUUID"] = sessionUUID.String()
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot generate token"})
	}

	expValue := claims["exp"].(int64)

	_, err = client.ActiveSessions.CreateOne(
		db.ActiveSessions.User.Link(
			db.User.ID.Equals(user.ID),
		),
		db.ActiveSessions.SessionUUID.Set(sessionUUID.String()),
		db.ActiveSessions.Exp.Set(time.Unix(expValue, 0)),
	).Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create session"})
	}

	setAuthCookie(c, t)

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func Login(c *fiber.Ctx, client *db.PrismaClient, jwtSecret []byte) error {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	ctx := context.Background()
	user, err := client.User.FindUnique(
		db.User.Email.Equals(input.Email),
	).Exec(ctx)
	if err == db.ErrNotFound {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	sessionUUID := uuid.New()

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.ID
	claims["sessionUUID"] = sessionUUID.String()
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString(jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot generate token"})
	}

	expValue := claims["exp"].(int64)

	_, err = client.ActiveSessions.CreateOne(
		db.ActiveSessions.User.Link(
			db.User.ID.Equals(user.ID),
		),
		db.ActiveSessions.SessionUUID.Set(sessionUUID.String()),
		db.ActiveSessions.Exp.Set(time.Unix(expValue, 0)),
	).Exec(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create session"})
	}

	setAuthCookie(c, t)

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}

func setAuthCookie(c *fiber.Ctx, token string) {
	cookie := new(fiber.Cookie)
	cookie.Name = "auth"
	cookie.Value = token
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.HTTPOnly = true
	cookie.Secure = false // Set to true in production
	cookie.SameSite = "Lax"
	c.Cookie(cookie)
}

func clearAuthCookie(c *fiber.Ctx) {
	c.ClearCookie("auth")
}

func Logout(c *fiber.Ctx, client *db.PrismaClient) error {
	ctx := context.Background()

	userID, ok := c.Locals("userID").(float64)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID in JWT"})
	}

	sessionUUID, ok := c.Locals("sessionUUID").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid session UUID in JWT"})
	}

	_, err := client.ActiveSessions.FindMany(
		db.ActiveSessions.UserID.Equals(int(userID)),
		db.ActiveSessions.SessionUUID.Equals(sessionUUID),
	).Delete().Exec(ctx)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{"error": "Failed to delete session"})
	}

	clearAuthCookie(c)

	// Close the SSE connection by aborting the context
	c.Context().Done()

	return c.SendStatus(fiber.StatusOK)
}

func GetUserProfile(c *fiber.Ctx, client *db.PrismaClient) error {
	userID := int(c.Locals("userID").(float64))

	ctx := context.Background()
	user, err := client.User.FindUnique(
		db.User.ID.Equals(userID),
	).Exec(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(fiber.Map{
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
		},
	})
}
