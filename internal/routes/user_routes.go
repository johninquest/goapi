// internal/routes/routes.go
package routes

import (
	"goapi/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
    // User route group
    user := app.Group("/api/user")
    user.Get("/", handlers.GetUser)
}