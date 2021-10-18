package ue

import (
	UEBlog "github.com/LittleBenx86/Benlog/internal/app/controller/ue"
	"github.com/LittleBenx86/Benlog/internal/app/middlewares"
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/LittleBenx86/Benlog/internal/global/variables"
	"github.com/gofiber/fiber/v2"
)

const UE_BLOG_V1_URL_PRE = "/ue/v1"

func UEV1BlogGroups(app *fiber.App) *fiber.App {
	blogHandler := UEBlog.NewBlogController(
		dependencies.WithDBClient(variables.DBInstance),
		dependencies.WithRedisClient(variables.RedisInstance),
		dependencies.WithLogger(variables.Logger))

	blogNoAuthGroup := app.Group(UE_BLOG_V1_URL_PRE+"/blog", middlewares.NoAuth)
	{
		blogNoAuthGroup.Get("/:blogId", blogHandler.Detail())
	}

	blogsNoAuthGroup := app.Group(UE_BLOG_V1_URL_PRE+"/blogs", middlewares.NoAuth)
	{
		blogsNoAuthGroup.Get("", blogHandler.List())
		blogsNoAuthGroup.Get("/page/:pageIndex", blogHandler.Page())
	}
	return app
}
