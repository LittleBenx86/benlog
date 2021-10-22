package middlewares

import (
	"errors"
	"github.com/LittleBenx86/Benlog/internal/app/model"
	"github.com/LittleBenx86/Benlog/internal/app/service"
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/LittleBenx86/Benlog/internal/global/variables"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func Anonymous(ctx *fiber.Ctx) error {
	// TODO anonymous token validation
	tokenService := service.NewSecurityService(false, nil)
	token := ctx.Get(fiber.HeaderAuthorization)
	i := strings.Index(token, "Bearer ")
	token = token[i+len("Bearer "):]
	_, err := tokenService.TokenBuilder.Parse(token) // If passed, the user must with anonymous role
	if err != nil {
		return err
	}

	// TODO anonymous security access validation
	a := model.Author{}
	a.UpdateRoleByAuthority()
	if ok, err := variables.SecurityEnforcer.Enforce(a,
		string(ctx.Request().URI().Path()),
		string(ctx.Request().Header.Method())); err != nil {
		return err
	} else {
		if !ok {
			return errors.New("access deny")
		}
	}

	return ctx.Next()
}

func Auth(ctx *fiber.Ctx) error {
	tokenService := service.NewSecurityService(true, &dependencies.Dependencies{
		DBClient: variables.DBClient,
	})
	token := ctx.Get(fiber.HeaderAuthorization)
	_, err := tokenService.TokenBuilder.Parse(token)
	if err != nil {
		return err
	}
	return ctx.Next()
}
