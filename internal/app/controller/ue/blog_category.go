package ue

import (
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/gofiber/fiber/v2"
)

type CategoryController struct {
	*dependencies.Dependencies
}

func NewCategoryController(options ...dependencies.Option) CategoryController {
	d := &dependencies.Dependencies{}

	for _, optionFn := range options {
		optionFn(d)
	}

	return CategoryController{
		Dependencies: d,
	}
}

func (c CategoryController) List() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
