package ue

import (
	"github.com/LittleBenx86/Benlog/internal/app/middlewares"
	"github.com/gofiber/fiber/v2"
)

func UEV1BlogTagGroups(app *fiber.App) *fiber.App {
	blogTagNoAuthGroup := app.Group(UE_BLOG_V1_URL_PRE+"/tags", middlewares.NoAuth)
	{
		blogTagNoAuthGroup.Get("")
	}
	return app
}
