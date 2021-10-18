package ue

import (
	"github.com/LittleBenx86/Benlog/internal/app/middlewares"
	"github.com/gofiber/fiber/v2"
)

func UEV1BlogCategoryGroups(app *fiber.App) *fiber.App {
	blogCategoryNoAuthGroup := app.Group(UE_BLOG_V1_URL_PRE+"/categories", middlewares.NoAuth)
	{
		blogCategoryNoAuthGroup.Get("")
	}
	return app
}
