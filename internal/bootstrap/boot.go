package bootstrap

import (
	"github.com/LittleBenx86/Benlog/internal/utils/convertor"
	"github.com/LittleBenx86/Benlog/internal/utils/websocket/ginws"
	"log"
	"os"

	InternalCasbin "github.com/LittleBenx86/Benlog/internal/app/security"
	"github.com/LittleBenx86/Benlog/internal/global/cfg"
	"github.com/LittleBenx86/Benlog/internal/global/consts"
	"github.com/LittleBenx86/Benlog/internal/global/variables"
	InternalContainer "github.com/LittleBenx86/Benlog/internal/repository/container"
	InternalDB "github.com/LittleBenx86/Benlog/internal/repository/mysql"
	InternalRedis "github.com/LittleBenx86/Benlog/internal/repository/redis"
	YmlFiles "github.com/LittleBenx86/Benlog/internal/utils/files"
	InternalLogger "github.com/LittleBenx86/Benlog/internal/utils/logger"
)

func init() {
	variables.BasePath = YmlFiles.GetProjectRuntimeRootPath() // init the application root path, i.e. the base path

	checkRequiredFolders()

	tmpInitYml := cfg.NewYmlConfigFactory(variables.BasePath)
	variables.YmlAppConfig = tmpInitYml.Clone("application-" + tmpInitYml.GetString("Configuration.Active.App"))
	variables.YmlAppConfig.ConfigFileUpdateListen()

	var zapCfg InternalLogger.ZapLoggerConfigParams
	if err := convertor.Map2Struct(variables.YmlAppConfig.Get("Logs").(map[string]interface{}), &zapCfg); err != nil {
		log.Fatal("unable to init logger with wrong yaml configuration")
	}
	InternalLogger.NewZapContext(zapCfg, variables.BasePath,
		variables.YmlAppConfig.GetString("App.Env"),
		InternalLogger.ZapLoggerHookHandler)
	variables.Logger = InternalLogger.GetInstance()

	var redisCfg InternalRedis.PoolConfiguration
	if err := convertor.Map2Struct(variables.YmlAppConfig.Get("Redis").(map[string]interface{}), &redisCfg); err != nil {
		log.Fatal("unable to init redis client with wrong yaml configuration")
	}
	InternalRedis.NewRedisContext(redisCfg)
	variables.RedisInstance = InternalRedis.GetInstance()

	if variables.YmlAppConfig.GetBool("Gormv2.MySQL.IsInitDBStartup") {
		var dbCfg InternalDB.DBConfigParams
		if err := convertor.Map2Struct(variables.YmlAppConfig.Get("Gormv2.MySQL").(map[string]interface{}), &dbCfg); err != nil {
			log.Fatal("unable to init database client with wrong yaml configuration")
		}
		InternalDB.NewDBContext(dbCfg, variables.Logger)
		variables.DBInstance = InternalDB.GetInstance()
	}

	if variables.YmlAppConfig.GetBool("Gormv2.MySQL.IsInitDBStartup") &&
		variables.YmlAppConfig.GetBool("Casbinv2.InitStartup") {
		InternalCasbin.Init("mysql", &variables.YmlAppConfig)
		variables.SecurityInstance = InternalCasbin.GetEnforcer()
	}

	InternalContainer.New(variables.Logger)

	if variables.YmlAppConfig.GetBool("Websocket.HubEnable") {
		variables.WebsocketHub = ginws.GetHubInstance()
		variables.WebsocketHub.Run()
	}
}

func checkRequiredFolders() {
	if _, err := os.Stat(variables.BasePath + "/conf/application.yml"); err != nil {
		log.Fatalf("%s, details: %s", consts.ERRORS_CONFIG_APP_YML_FILE_NOT_EXISTS, err.Error())
	}
}
