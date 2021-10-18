package routers

import (
	"github.com/gofiber/fiber/v2"
)

func InitAEHttpServer() *fiber.App {
	adminApp := fiber.New()
	return adminApp
}
