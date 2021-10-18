package ue

import (
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/gofiber/fiber/v2"
)

type CommentController struct {
	*dependencies.Dependencies
}

func NewCommentController(options ...dependencies.Option) CommentController {
	d := &dependencies.Dependencies{}

	for _, optionFn := range options {
		optionFn(d)
	}

	return CommentController{
		Dependencies: d,
	}
}

func (c CommentController) Create() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}

func (c CommentController) Page() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		return nil
	}
}
