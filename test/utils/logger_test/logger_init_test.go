package logger_test

import (
	"github.com/LittleBenx86/Benlog/internal/global/cfg"
	"github.com/LittleBenx86/Benlog/internal/utils/convertor"
	InternalLogger "github.com/LittleBenx86/Benlog/internal/utils/logger"
	"testing"
)

func Test_NewZapContext(t *testing.T) {
	rootPath := "/Users/benzheng/Projects/Go/Benlog"
	tmpInitYml := cfg.NewYmlConfigFactory(rootPath)
	ymlAppConfig := tmpInitYml.Clone("application-" + tmpInitYml.GetString("Configuration.Active.App"))

	var zapCfg InternalLogger.ZapLoggerConfigParams
	t.Logf("%+v\n", ymlAppConfig.Get("Logs").(map[string]interface{}))
	if err := convertor.Map2Struct(ymlAppConfig.Get("Logs").(map[string]interface{}), &zapCfg); err != nil {
		t.Fatal("unable to init logger with wrong yaml config")
	}
	t.Logf("%+v\n", zapCfg)
}
