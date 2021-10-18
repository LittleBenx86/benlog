package ue

import (
	"errors"
	"github.com/LittleBenx86/Benlog/internal/app/response"
	UEService "github.com/LittleBenx86/Benlog/internal/app/service/ue"
	"github.com/LittleBenx86/Benlog/internal/global/consts"
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type MetricsController struct {
	*dependencies.Dependencies
	MetricsType MetricsType
}

type MetricsType int

const (
	METRICS_ALL MetricsType = 2*iota + 1
	METRICS_APP_MEM
	METRICS_APP_CPU
	METRICS_APP_HEAP
	METRICS_APP_MUTEX
	METRICS_APP_THREADS
)

func NewMetricsController(metricsType MetricsType, options ...dependencies.Option) MetricsController {
	d := &dependencies.Dependencies{}

	for _, optionFn := range options {
		optionFn(d)
	}

	return MetricsController{
		MetricsType:  metricsType,
		Dependencies: d,
	}
}

func (m MetricsController) Clone(metricsType MetricsType) *MetricsController {
	return &MetricsController{
		MetricsType:  metricsType,
		Dependencies: m.Dependencies,
	}
}

func (m MetricsController) Detail() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		svc := UEService.NewMetricsService(m.Dependencies)
		var out string
		var err error
		switch m.MetricsType {
		default:
			out = ""
			err = errors.New("unknown metrics type")
		case METRICS_ALL:
			out, err = svc.GetMetrics()
			goto responseHandle
		case METRICS_APP_CPU:
			out, err = svc.GetCpuMetrics()
			goto responseHandle
		case METRICS_APP_MEM:
			out, err = svc.GetMemMetrics()
			goto responseHandle
		}

	responseHandle:
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
				SetResponseJson(out).
				Ok()
		}
		return responseErr
	}
}
