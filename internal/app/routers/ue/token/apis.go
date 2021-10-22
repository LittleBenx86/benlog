package token

import (
	"github.com/LittleBenx86/Benlog/internal/app/controller/ue"
	"github.com/gofiber/fiber/v2"
)

const UE_TOKEN_V1_URL_PRE = "/ue/v1/token"

func UEV1TokenGroups(app *fiber.App) *fiber.App {

	tokenController := ue.NewTokenController()
	metricsNoAuthGroup := app.Group(UE_TOKEN_V1_URL_PRE)
	{
		metricsNoAuthGroup.Get("", tokenController.Detail())
	}

	return app
}
