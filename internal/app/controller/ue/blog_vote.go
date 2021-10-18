package ue

import (
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/gofiber/fiber/v2"
)

type VoteController struct {
	*dependencies.Dependencies
}

func NewVoteController(options ...dependencies.Option) VoteController {
	d := &dependencies.Dependencies{}

	for _, optionFn := range options {
		optionFn(d)
	}

	return VoteController{
		Dependencies: d,
	}
}

func (v VoteController) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func (v VoteController) Delete() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func (v VoteController) Detail() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
