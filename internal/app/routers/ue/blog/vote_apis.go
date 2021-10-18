package ue

import (
	"github.com/LittleBenx86/Benlog/internal/app/middlewares"
	"github.com/gofiber/fiber/v2"
)

func UEV1BlogVoteGroups(app *fiber.App) *fiber.App {
	blogVoteNoAuthGroup := app.Group(UE_BLOG_V1_URL_PRE+"/vote", middlewares.NoAuth)
	{
		blogVoteNoAuthGroup.Get("/:blogId")
		blogVoteNoAuthGroup.Post("/:blogId")
		blogVoteNoAuthGroup.Delete("/:blogId")
	}
	return app
}
