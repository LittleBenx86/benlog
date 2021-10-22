package metrics

import (
	"github.com/LittleBenx86/Benlog/internal/app/controller/ue"
	"github.com/LittleBenx86/Benlog/internal/app/middlewares"
	"github.com/LittleBenx86/Benlog/internal/global/dependencies"
	"github.com/LittleBenx86/Benlog/internal/global/variables"
	"github.com/gofiber/fiber/v2"
)

const UE_Metrics_V1_URL_PRE = "/ue/v1/metrics"

func UEV1MetricsGroups(app *fiber.App) *fiber.App {

	metricsAllController := ue.NewMetricsController(ue.METRICS_ALL, dependencies.WithLogger(variables.Logger))
	metricsNoAuthGroup := app.Group(UE_Metrics_V1_URL_PRE, middlewares.Anonymous)
	{
		metricsNoAuthGroup.Get("", metricsAllController.Detail())
		metricsNoAuthGroup.Get("/cpu", metricsAllController.Clone(ue.METRICS_APP_CPU).Detail())
		metricsNoAuthGroup.Get("/mem", metricsAllController.Clone(ue.METRICS_APP_MEM).Detail())
	}

	return app
}
