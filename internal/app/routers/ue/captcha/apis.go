package captcha

import (
	"github.com/LittleBenx86/Benlog/internal/app/middlewares"
	"github.com/gofiber/fiber/v2"
)

const UE_RESUME_V1_URL_PRE = "/ue/author/info"

func UEV1ResumeGroups(app *fiber.App) *fiber.App {
	resumeNoAuthGroup := app.Group(UE_RESUME_V1_URL_PRE, middlewares.NoAuth)
	{
		resumeNoAuthGroup.Post("")
	}
	return app
}
