package admin

import (
	"github.com/LittleBenx86/Benlog/internal/app/middlewares"
	"github.com/gofiber/fiber/v2"
)

const ADMIN_V1_BLOG_URL_PRE = "/admin/v1/blog"

func AdminV1BlogGroups(app *fiber.App) *fiber.App {
	apiGroup := app.Group(ADMIN_V1_BLOG_URL_PRE, middlewares.Anonymous)
	{
		apiGroup.Get("/:blogId")
		apiGroup.Post("/:blogId")
		apiGroup.Delete("/:blogId")
		apiGroup.Patch("/:blogId")
	}
	return app
}
