package admin

import (
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/gofiber/fiber/v2"
)

type BlogController struct {
	*dependencies.Dependencies
}

func NewBlogController(options ...dependencies.Option) BlogController {
	d := &dependencies.Dependencies{}

	for _, optionFn := range options {
		optionFn(d)
	}

	return BlogController{
		Dependencies: d,
	}
}

func (b *BlogController) Detail() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func (b *BlogController) Page() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
