package ue

import (
	"github.com/LittleBenx86/Benlog/internal/app/middlewares"
	"github.com/gofiber/fiber/v2"
)

func UEV1BlogViewGroups(app *fiber.App) *fiber.App {
	blogViewNoAuthGroup := app.Group(UE_BLOG_V1_URL_PRE+"/view", middlewares.NoAuth)
	{
		blogViewNoAuthGroup.Post("/:blogId")
		blogViewNoAuthGroup.Delete("/:blogId")
		blogViewNoAuthGroup.Get("/:blogId")
	}
	return app
}
