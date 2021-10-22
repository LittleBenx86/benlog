package ue

import (
	"github.com/LittleBenx86/Benlog/internal/app/middlewares"
	"github.com/gofiber/fiber/v2"
)

const UE_COMM_V1_URL_PRE = "/ue/v1/comm"

func UEV1CommentGroups(app *fiber.App) *fiber.App {
	apiGroup := app.Group(UE_COMM_V1_URL_PRE, middlewares.Anonymous)
	{
		apiGroup.Get("/:blogId")
		apiGroup.Post("/:blogId")
	}
	return app
}
