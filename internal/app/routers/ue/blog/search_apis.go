package ue

import (
	"github.com/LittleBenx86/Benlog/internal/app/middlewares"
	"github.com/gofiber/fiber/v2"
)

func UEV1BlogSearchGroups(app *fiber.App) *fiber.App {
	blogSearchNoAuthGroup := app.Group(UE_BLOG_V1_URL_PRE+"/search", middlewares.NoAuth)
	{
		blogSearchNoAuthGroup.Post("")
	}
	return app
}
