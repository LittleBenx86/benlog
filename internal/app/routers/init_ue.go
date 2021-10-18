package routers

import (
	"log"
	"os"
	"time"

	"github.com/LittleBenx86/Benlog/internal/app/middlewares"
	"github.com/LittleBenx86/Benlog/internal/app/response"
	"github.com/LittleBenx86/Benlog/internal/app/routers/ue/metrics"
	"github.com/LittleBenx86/Benlog/internal/global/variables"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func InitUEHttpServer() *fiber.App {
	userApp := fiber.New(fiber.Config{
		ErrorHandler: response.DefaultErrorHandler,
	})

	appMode := variables.YmlAppConfig.GetString("App.Env")
	if appMode == "dev" {
		userApp.Use(pprof.New())

		userApp.Use(logger.New(logger.Config{
			Output:       os.Stdout,
			Format:       "[${time}] <${blue}${pid}${reset}> ${status} - ${latency} ${method} ${path}\n",
			TimeFormat:   "2006/01/02 15:04:05.000",
			TimeZone:     "Asia/Shanghai",
			TimeInterval: 500 * time.Millisecond,
			Next:         nil,
		}))
	} else if appMode == "prod" {
		// without nginx proxy, disable output log to console
		f, _ := os.OpenFile(variables.BasePath+variables.YmlAppConfig.GetString("Logs.GinLogOutput"),
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				variables.Logger.Error("close fiber log file error: " + err.Error())
			}
		}(f)

		userApp.Use(logger.New(logger.Config{
			Output:       f,
			Format:       "[${time}] <${pid}> ${status} - ${latency} ${method} ${path}\n",
			TimeFormat:   "2006/01/02 15:04:05.000",
			TimeZone:     "Asia/Shanghai",
			TimeInterval: 500 * time.Millisecond,
			Next:         nil,
		}))
	} else {
		log.Fatalf("unknown app env mode [%s] to boot server\n", appMode)
	}

	// cors
	if variables.YmlAppConfig.GetBool("App.Server.AllowCrossDomain") {
		userApp.Use(middlewares.CORSNext())
	}

	userApp = metrics.UEV1MetricsGroups(userApp)

	return userApp
}
