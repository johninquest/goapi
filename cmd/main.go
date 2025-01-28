package main

import (
	"fmt"
	"log"
	"os"

	"goapi/internal/database"
	"goapi/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize SQLite database
	db, err := database.NewSQLiteDB()
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Error running migrations:", err)
	}

	// Initialize Fiber
	app := fiber.New()
	app.Use(logger.New())

	// Setup routes
	routes.UserRoutes(app) 
	routes.HelloRoutes(app)
	routes.PolicyRoutes(app, db)

	// Public routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello World from Go/Fiber!",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
