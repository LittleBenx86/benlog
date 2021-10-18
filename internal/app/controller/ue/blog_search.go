package ue

import (
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/gofiber/fiber/v2"
)

type SearchController struct {
	*dependencies.Dependencies
}

func NewSearchController(options ...dependencies.Option) SearchController {
	d := &dependencies.Dependencies{}

	for _, optionFn := range options {
		optionFn(d)
	}

	return SearchController{
		Dependencies: d,
	}
}

func (s SearchController) Page() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
