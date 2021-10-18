package mysql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

func NewCustomGormDBLog(dbType string, options ...Options) gormLogger.Interface {
	var (
		infoFmt       = "[info] %s %s \n"
		warnFmt       = "[warn] %s %s \n"
		errorFmt      = "[error] %s %s \n"
		traceFmt      = "[trace] %s %s [%.3fms] [rows:%v] %s \n"
		traceWarnFmt  = "[traceWarn] %s %s %s [%.3fms] [rows:%v] %s \n"
		traceErrorFmt = "[traceError] %s %s %s [%.3fms] [rows:%v] %s \n"
	)

	logConf := gormLogger.Config{
		SlowThreshold: time.Second * dbContext.Cfg.GlobalSlowSQLThreshold,
		LogLevel:      gormLogger.Warn,
		Colorful:      false,
	}

	log := &DBLogger{
		Writer:        logOutput{},
		Config:        logConf,
		infoFmt:       infoFmt,
		warnFmt:       warnFmt,
		errorFmt:      errorFmt,
		traceFmt:      traceFmt,
		traceWarnFmt:  traceWarnFmt,
		traceErrorFmt: traceErrorFmt,
	}

	for _, val := range options {
		val.apply(log)
	}
	return log
}

type DBLogger struct {
	gormLogger.Writer
	gormLogger.Config
	infoFmt, warnFmt, errorFmt            string
	traceFmt, traceErrorFmt, traceWarnFmt string
}

func (l *DBLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	loggerDeepCopy := *l
	loggerDeepCopy.LogLevel = level
	return &loggerDeepCopy
}

func (l *DBLogger) Info(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Info {
		l.Printf(l.infoFmt, append([]interface{}{time.Now(), msg}, data...)...)
	}
}

func (l *DBLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Warn {
		l.Printf(l.warnFmt, append([]interface{}{time.Now(), msg}, data...)...)
	}
}

func (l *DBLogger) Error(_ context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormLogger.Error {
		l.Printf(l.errorFmt, append([]interface{}{time.Now(), msg}, data...)...)
	}
}

func (l *DBLogger) Trace(_ context.Context, begin time.Time, fn func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		case err != nil && l.LogLevel >= gormLogger.Error:
			sql, rows := fn()
			if rows == -1 {
				l.Printf(l.traceErrorFmt, time.Now(), utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-1", sql)
			} else {
				l.Printf(l.traceErrorFmt, time.Now(), utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= gormLogger.Warn:
			sql, rows := fn()
			slowLog := fmt.Sprintf("Slow SQL >= %v", l.SlowThreshold)
			if rows == -1 {
				l.Printf(l.traceWarnFmt, time.Now(), utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-1", sql)
			} else {
				l.Printf(l.traceWarnFmt, time.Now(), utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		case l.LogLevel >= gormLogger.Info:
			sql, rows := fn()
			if rows == -1 {
				l.Printf(l.traceFmt, time.Now(), utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-1", sql)
			} else {
				l.Printf(l.traceFmt, time.Now(), utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
			}
		}
	}
}

type Options interface {
	apply(*DBLogger)
}

func SetInfoFormat(format string) Options {
	return OptionsFn(func(l *DBLogger) {
		l.infoFmt = format
	})
}

func SetWarnFormat(format string) Options {
	return OptionsFn(func(l *DBLogger) {
		l.warnFmt = format
	})
}

func SetErrorFormat(format string) Options {
	return OptionsFn(func(l *DBLogger) {
		l.errorFmt = format
	})
}

func SetTraceFormat(format string) Options {
	return OptionsFn(func(l *DBLogger) {
		l.traceFmt = format
	})
}

func SetTraceWarnFormat(format string) Options {
	return OptionsFn(func(l *DBLogger) {
		l.traceWarnFmt = format
	})
}

func SetTraceErrorFormat(format string) Options {
	return OptionsFn(func(l *DBLogger) {
		l.traceErrorFmt = format
	})
}

type OptionsFn func(logger *DBLogger)

func (fn OptionsFn) apply(l *DBLogger) {
	fn(l)
}

type logOutput struct{}

func (l logOutput) Printf(strFmt string, args ...interface{}) {
	logResult := fmt.Sprintf(strFmt, args...)
	logFlag := "gorm-v2 mysql logs: "
	detailFlag := "details: "
	if strings.Contains(strFmt, "[info]") || strings.Contains(strFmt, "[trace]") {
		dbContext.Logger.Info(logFlag, zap.String(detailFlag, logResult))
	} else if strings.Contains(strFmt, "[error]") || strings.Contains(strFmt, "[traceError]") {
		dbContext.Logger.Error(logFlag, zap.String(detailFlag, logResult))
	} else if strings.Contains(strFmt, "[warn]") || strings.Contains(strFmt, "[traceWarn]") {
		dbContext.Logger.Warn(logFlag, zap.String(detailFlag, logResult))
	}
}
