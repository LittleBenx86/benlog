package ue

import (
	"github.com/LittleBenx86/Benlog/internal/app/response"
	"github.com/LittleBenx86/Benlog/internal/app/service"
	"github.com/LittleBenx86/Benlog/internal/global/consts"
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/LittleBenx86/Benlog/internal/utils/uuid"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type TokenController struct {
	*dependencies.Dependencies
}

func NewTokenController(options ...dependencies.Option) TokenController {
	d := &dependencies.Dependencies{}

	for _, optionFn := range options {
		optionFn(d)
	}

	return TokenController{
		Dependencies: d,
	}
}

func (t TokenController) Detail() fiber.Handler {
	tokenService := service.NewSecurityService(false, t.Dependencies)
	return func(ctx *fiber.Ctx) error {
		var uid, token string
		var err error

		if uid, err = uuid.GenerateRandomStringId(16); err != nil {
			goto responseErrorHandle
		}

		if token, err = tokenService.Generate(service.CustomClaims{
			UID:  uid,
			Name: consts.USER_ANONYMOUS_NAME,
		}); err != nil {
			goto responseErrorHandle
		}

	responseErrorHandle:
		var responseErr error
		if err != nil {
			responseErr = response.NewStream(ctx).
				SetHttpCode(http.StatusInternalServerError).
				SetAppCode(consts.AppCommonInternalError).
				SetDetails(err.Error()).
				Fail()
		} else {
			responseErr = response.NewStream(ctx).
				SetHttpCode(http.StatusOK).
				SetResponseJson(token).
				Ok()
		}
		return responseErr
	}
}

func (t TokenController) Modify() fiber.Handler {
	panic("implement me")
}
