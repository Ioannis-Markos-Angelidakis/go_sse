package main

import (
	"backend/broker"
	"backend/prisma/db"
	"backend/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

var client = db.NewClient()
var jwtSecret []byte
var taskBroker = broker.NewBroker() // Changed variable name

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Load JWT secret from environment
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	if len(jwtSecret) == 0 {
		panic("JWT_SECRET environment variable not set")
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://192.168.0.13:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	if err := client.Prisma.Connect(); err != nil {
		panic(err)
	}
	defer client.Prisma.Disconnect()

	routes.SetupRoutes(app, client, taskBroker, jwtSecret) // Updated variable name

	port := os.Getenv("BACKEND_PORT")
	if port == "" {
		port = "5000"
	}
	app.Listen("192.168.0.13:" + port)
}
