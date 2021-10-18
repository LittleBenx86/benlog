package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSNext() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Access-Control-Allow-Headers, Authorization,User-Agent, Keep-Alive, Content-Type, X-Requested-With, X-CSRF-Token, AccessToken, Token",
		ExposeHeaders:    "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type",
		AllowMethods:     "GET, POST, DELETE, PUT, PATCH",
		AllowCredentials: true,
		Next:             nil,
	})
}
