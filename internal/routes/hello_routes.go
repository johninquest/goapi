// internal/routes/routes.go
package routes

import (
	"goapi/internal/handlers"

	"github.com/gofiber/fiber/v2"
)

func HelloRoutes(app *fiber.App) {
	// User route group
	user := app.Group("/api/hello")
	user.Get("/", handlers.GetHello)
}
