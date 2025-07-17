package routes

import (
	"backend/broker"
	"backend/middleware"
	"backend/prisma/db"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, client *db.PrismaClient, broker *broker.Broker, jwtSecret []byte) {
	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	app.Post("/register", func(c *fiber.Ctx) error {
		return Register(c, client, jwtSecret)
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		return Login(c, client, jwtSecret)
	})
	app.Get("/me", authMiddleware, func(c *fiber.Ctx) error {
		return GetUserProfile(c, client)
	})
	app.Post("/logout", authMiddleware, func(c *fiber.Ctx) error {
		return Logout(c, client)
	})
	app.Get("/events", authMiddleware, func(c *fiber.Ctx) error {
		return SSEHandler(c, broker)
	})
	app.Get("/public-tasks", func(c *fiber.Ctx) error {
		return GetPublicTasks(c, client)
	})

	taskGroup := app.Group("/tasks")
	taskGroup.Use(authMiddleware)
	taskGroup.Post("/", func(c *fiber.Ctx) error {
		return CreateTask(c, client, broker)
	})
	taskGroup.Get("/", func(c *fiber.Ctx) error {
		return GetTasks(c, client)
	})
	taskGroup.Put("/:id", func(c *fiber.Ctx) error {
		return UpdateTask(c, client, broker)
	})
	taskGroup.Delete("/:id", func(c *fiber.Ctx) error {
		return DeleteTask(c, client, broker)
	})
}
