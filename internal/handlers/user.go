// internal/handlers/user.go
package handlers

import "github.com/gofiber/fiber/v2"

func GetUser(c *fiber.Ctx) error {
    return c.JSON(fiber.Map{
        "message": "Hello User",
    })
}