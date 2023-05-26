package auth

import "github.com/gofiber/fiber/v2"

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement auth handling with Saqura Auth Service
		return c.Next()
	}
}
