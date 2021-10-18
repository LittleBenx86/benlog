package middlewares

import "github.com/gofiber/fiber/v2"

func NoAuth(ctx *fiber.Ctx) error {
	return ctx.Next()
}
