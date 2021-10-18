package ue

import (
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/gofiber/fiber/v2"
)

type ViewController struct {
	*dependencies.Dependencies
}

func NewViewController(options ...dependencies.Option) ViewController {
	d := &dependencies.Dependencies{}

	for _, optionFn := range options {
		optionFn(d)
	}

	return ViewController{
		Dependencies: d,
	}
}

func (v ViewController) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func (v ViewController) Delete() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func (v ViewController) Detail() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
