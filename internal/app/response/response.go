package response

import (
	"net/http"

	"github.com/LittleBenx86/Benlog/internal/global/consts"
	InternalLogger "github.com/LittleBenx86/Benlog/internal/utils/logger"

	"github.com/gofiber/fiber/v2"
)

type coreResponse struct {
	appCode  consts.AppStatusCode
	httpCode int
	details  string
	data     interface{}
}

type appResponse struct {
	coreResponse
	requestUrl string
}

var DefaultErrorHandler = func(ctx *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	ctx.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)

	return ctx.Status(code).SendString(err.Error())
}

func (e *appResponse) ErrorDetails() string {
	return "From request url: [" + e.requestUrl +
		"]\n Internal error code: [" + string(rune(e.appCode)) +
		"]\n Error Details: " + e.details
}

type Stream struct {
	Context         *fiber.Ctx
	ResponseContent *appResponse
	JsonStrContent  string
	IsUseJsonStr    bool
}

func NewStream(ctx *fiber.Ctx) *Stream {
	return &Stream{
		Context: ctx,
		ResponseContent: &appResponse{
			coreResponse: coreResponse{
				httpCode: http.StatusOK,                 // default ok
				appCode:  consts.RequestCommonSucceeded, // default succeeded
				details:  "",
				data:     nil,
			},
			requestUrl: "",
		},
		JsonStrContent: "",
		IsUseJsonStr:   false,
	}
}

func (s *Stream) SetHttpCode(code int) *Stream {
	s.ResponseContent.httpCode = code
	return s
}

func (s *Stream) SetResponseJson(json string) *Stream {
	s.IsUseJsonStr = true
	s.JsonStrContent = json
	return s
}

func (s *Stream) SetAppCode(code consts.AppStatusCode) *Stream {
	s.ResponseContent.appCode = code
	return s
}

func (s *Stream) SetDetails(details string) *Stream {
	s.ResponseContent.details = details
	return s
}

func (s *Stream) SetAdditionalData(data interface{}) *Stream {
	s.ResponseContent.data = data
	return s
}

func (s *Stream) SetRequestUrl(url string) *Stream {
	s.ResponseContent.requestUrl = url
	return s
}

func (s *Stream) Ok() error {
	if s.IsUseJsonStr {
		return responseJsonFromStr(s.Context, s.ResponseContent.httpCode, s.JsonStrContent)
	}
	return responseJson(s.Context, s.ResponseContent)
}

func (s *Stream) Fail() error {
	var err error
	if s.IsUseJsonStr {
		err = responseJsonFromStr(s.Context, s.ResponseContent.httpCode, s.JsonStrContent)
	} else {
		err = responseJson(s.Context, s.ResponseContent)
	}
	InternalLogger.GetInstance().Error(s.ResponseContent.ErrorDetails())
	return err
}

func responseJson(ctx *fiber.Ctx, body *appResponse) error {
	return ctx.Status(body.httpCode).JSON(fiber.Map{
		"code": body.appCode,
		"msg":  body.details,
		"data": body.data,
	})
}

func responseJsonFromStr(ctx *fiber.Ctx, httpCode int, json string) error {
	ctx.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSONCharsetUTF8)
	return ctx.Status(httpCode).SendString(json)
}
