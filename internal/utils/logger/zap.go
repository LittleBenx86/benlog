package logger

import (
	"log"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLoggerConfigParams struct {
	GinLogOutputPath     string `json:"GinLogOutput" mapstructure:"ginlogoutput"`
	AppLogOutputPath     string `json:"AppLogOutput" mapstructure:"applogoutput"`
	OutputFormat         string `json:"OutputFormat" mapstructure:"outputformat"`
	TimePrecision        string `json:"TimePrecision" mapstructure:"timeprecision"`
	LogFileMaxSize       int    `json:"MaxSize" mapstructure:"maxsize"`
	LogFileMaxBackupDays int    `json:"MaxBackups" mapstructure:"maxbackups"`
	LogFileMaxLiveDays   int    `json:"MaxAge" mapstructure:"maxage"`
	CompressEnable       bool   `json:"Compress" mapstructure:"compress"`
}

type ZapContext struct {
	Env             string
	Cfg             *ZapLoggerConfigParams
	ProjectRootPath string
	EntryHookFn     func(entry zapcore.Entry) error
}

var (
	instance   *zap.Logger
	once       sync.Once
	zapContext *ZapContext
)

func NewZapContext(cfg ZapLoggerConfigParams, root string, env string, fn func(entry zapcore.Entry) error) {
	if fn == nil || len(root) == 0 || len(env) == 0 {
		log.Fatal("init zap logger with empty parameter!")
	}

	zapContext = &ZapContext{
		Env:             env,
		Cfg:             &cfg,
		ProjectRootPath: root,
		EntryHookFn:     fn,
	}
}

func GetInstance() *zap.Logger {
	once.Do(newZapLogger)
	return instance
}

func newZapLogger() {
	appMode := zapContext.Env

	if appMode == "dev" {
		// "dev" means debug enable, just generate a development zap log manager address, let logs all output into console
		if logger, err := zap.NewDevelopment(zap.Hooks(zapContext.EntryHookFn)); err == nil {
			instance = logger
			return
		} else {
			log.Fatalf("%s, details; %s", "Create zap log package failed", err.Error())
		}
	}

	if appMode != "prod" {
		log.Fatalf("Unknown app mode: %s", appMode)
	}

	// "prod" work mode configuration
	encoderConfig := zap.NewProductionEncoderConfig()
	var recordTimeFormat string
	switch zapContext.Cfg.TimePrecision {
	case "second":
		recordTimeFormat = "2006-01-02 15:04:05"
	case "millisecond":
		recordTimeFormat = "2006-01-02 15:04:05.000"
	default:
		recordTimeFormat = "2006/01/02 15:04:05"
	}

	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(recordTimeFormat))
	}

	encoderConfig.TimeKey = "created_at" // The json format datetime key is easy to export to ELK

	var encoder zapcore.Encoder
	switch zapContext.Cfg.OutputFormat {
	case "plaintext":
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case "json":
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		log.Fatalf("Unsupported output format type [%s]", zapContext.Cfg.OutputFormat)
	}

	logFileQualifiedName := zapContext.ProjectRootPath + zapContext.Cfg.AppLogOutputPath
	lumberLogger := &lumberjack.Logger{
		Filename:   logFileQualifiedName,
		MaxSize:    zapContext.Cfg.LogFileMaxSize,
		MaxBackups: zapContext.Cfg.LogFileMaxBackupDays,
		MaxAge:     zapContext.Cfg.LogFileMaxLiveDays,
		Compress:   zapContext.Cfg.CompressEnable,
	}
	writer := zapcore.AddSync(lumberLogger)
	zapCore := zapcore.NewCore(encoder, writer, zap.InfoLevel) // "prod" mode default log level is info
	instance = zap.New(zapCore, zap.AddCaller(), zap.Hooks(zapContext.EntryHookFn))
}
