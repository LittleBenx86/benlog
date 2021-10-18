package ue

import (
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/gofiber/fiber/v2"
)

type TagController struct {
	*dependencies.Dependencies
}

func NewTagController(options ...dependencies.Option) TagController {
	d := &dependencies.Dependencies{}

	for _, optionFn := range options {
		optionFn(d)
	}

	return TagController{
		Dependencies: d,
	}
}

func (t TagController) List() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
