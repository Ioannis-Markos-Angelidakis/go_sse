package main

import (
	"backend/broker"
	"backend/prisma/db"
	"backend/routes"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

var client = db.NewClient()
var jwtSecret []byte
var taskBroker = broker.NewBroker()

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		panic("Error loading BACKEND_PORT from .env")
	}

	host := os.Getenv("HOST")
	if host == "" {
		panic("Error loading HOST from .env")
	}

	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		panic("JWT_SECRET environment variable not set")
	}

	app := fiber.New()

	app.Use(
		recover.New(),
		logger.New(logger.Config{
			Format: "[${ip}]:${port} ${status} - ${method} ${path}\n",
		}),
		limiter.New(limiter.Config{
			Max:               20,
			Expiration:        30 * time.Second,
			LimiterMiddleware: limiter.SlidingWindow{},
		}),
		cors.New(cors.Config{
			AllowOrigins:     fmt.Sprintf("http://%s", host),
			AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
			AllowCredentials: true,
		}))

	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}
	defer client.Prisma.Disconnect()

	routes.SetupRoutes(app, client, taskBroker, jwtSecret)

	app.Listen(fmt.Sprintf("%s:%s", host, port))
}
