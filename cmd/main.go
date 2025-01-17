package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"goapi/internal/routes"
)

// Custom middleware - like Express middleware
/* func authMiddleware(c *fiber.Ctx) error {
	// Get token from header
	token := c.Get("Authorization")

	if token == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Continue to next middleware/handler
	return c.Next()
} */

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // default port if not specified
	}

	// Initialize Fiber
	app := fiber.New()

	// Setup routes
	routes.UserRoutes(app)

	// Global middleware - like app.use() in Express
	app.Use(logger.New())

	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello World from Go/Fiber!",
		})
	})

	/* // Protected routes group
	api := app.Group("/api", authMiddleware) // Apply middleware to group

	api.Get("/protected", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Protected route",
		})
	}) */

	// Start server
	// log.Fatal(app.Listen(":3000"))
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
