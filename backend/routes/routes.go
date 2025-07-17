package routes

import (
	"backend/broker"
	"backend/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, broker *broker.Broker, jwtSecret []byte) {
	authMiddleware := middleware.AuthMiddleware(jwtSecret)

	app.Post("/register", func(c *fiber.Ctx) error {
		return Register(c, jwtSecret)
	})
	app.Post("/login", func(c *fiber.Ctx) error {
		return Login(c, jwtSecret)
	})
	app.Get("/me", authMiddleware, func(c *fiber.Ctx) error {
		return GetUserProfile(c)
	})
	app.Post("/logout", authMiddleware, func(c *fiber.Ctx) error {
		return Logout(c)
	})
	app.Get("/events", authMiddleware, func(c *fiber.Ctx) error {
		return SSEHandler(c, broker)
	})
	app.Get("/public-tasks", func(c *fiber.Ctx) error {
		return GetPublicTasks(c)
	})

	taskGroup := app.Group("/tasks")
	taskGroup.Use(authMiddleware)
	taskGroup.Post("/", func(c *fiber.Ctx) error {
		return CreateTask(c, broker)
	})
	taskGroup.Get("/", func(c *fiber.Ctx) error {
		return GetTasks(c)
	})
	taskGroup.Put("/:id", func(c *fiber.Ctx) error {
		return UpdateTask(c, broker)
	})
	taskGroup.Delete("/:id", func(c *fiber.Ctx) error {
		return DeleteTask(c, broker)
	})
}
