package ue

import (
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/gofiber/fiber/v2"
)

type CaptchaController struct {
	*dependencies.Dependencies
}

func NewCaptchaController(options ...dependencies.Option) CaptchaController {
	d := &dependencies.Dependencies{}

	for _, optionFn := range options {
		optionFn(d)
	}

	return CaptchaController{
		Dependencies: d,
	}
}

func (r CaptchaController) Detail() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
