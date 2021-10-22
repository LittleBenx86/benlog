package bootstrap

import (
	"log"
	"os"

	"github.com/LittleBenx86/Benlog/internal/app/model"
	"github.com/LittleBenx86/Benlog/internal/global/cfg"
	"github.com/LittleBenx86/Benlog/internal/global/consts"
	"github.com/LittleBenx86/Benlog/internal/global/variables"
	InternalContainer "github.com/LittleBenx86/Benlog/internal/repository/container"
	InternalDB "github.com/LittleBenx86/Benlog/internal/repository/mysql"
	InternalRedis "github.com/LittleBenx86/Benlog/internal/repository/redis"
	"github.com/LittleBenx86/Benlog/internal/utils/convertor"
	YmlFiles "github.com/LittleBenx86/Benlog/internal/utils/files"
	InternalLogger "github.com/LittleBenx86/Benlog/internal/utils/logger"
	InternalCasbin "github.com/LittleBenx86/Benlog/internal/utils/security"
	"github.com/LittleBenx86/Benlog/internal/utils/websocket/fiberws"
)

func init() {
	variables.BasePath = YmlFiles.GetProjectRuntimeRootPath() // init the application root path, i.e. the base path

	checkRequiredFolders()

	// TODO load yml configurations
	tmpInitYml := cfg.NewYmlConfigFactory(variables.BasePath)
	variables.YmlAppConfig = tmpInitYml.Clone("application-" + tmpInitYml.GetString("Configuration.Active.App"))
	variables.YmlAppConfig.ConfigFileUpdateListen()

	// TODO load zap logger
	var zapCfg InternalLogger.ZapLoggerConfigParams
	if err := convertor.Map2Struct(variables.YmlAppConfig.Get("Logs").(map[string]interface{}), &zapCfg); err != nil {
		log.Fatal("unable to init logger with wrong yaml configuration")
	}
	InternalLogger.NewZapContext(zapCfg, variables.BasePath,
		variables.YmlAppConfig.GetString("App.Env"),
		InternalLogger.ZapLoggerHookHandler)
	variables.Logger = InternalLogger.GetInstance()

	// TODO load redis client
	var redisCfg InternalRedis.PoolConfiguration
	if err := convertor.Map2Struct(variables.YmlAppConfig.Get("Redis").(map[string]interface{}), &redisCfg); err != nil {
		log.Fatal("unable to init redis client with wrong yaml configuration")
	}
	InternalRedis.NewRedisContext(redisCfg)
	variables.RedisClient = InternalRedis.GetInstance()

	// TODO load database client
	var dbCfg InternalDB.DBConfigParams
	if err := convertor.Map2Struct(variables.YmlAppConfig.Get("Gormv2.MySQL").(map[string]interface{}), &dbCfg); err != nil {
		log.Fatal("unable to init database client with wrong yaml configuration")
	}
	InternalDB.NewDBContext(dbCfg, variables.Logger)
	variables.DBClient = InternalDB.GetInstance()
	variables.DBClient.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8")

	// TODO databases migration
	if dbCfg.InitDBStartup {
		migrateRequiredTables()
	}

	// TODO load casbin securities
	var casbinCfg InternalCasbin.CasbinConfiguration
	if variables.YmlAppConfig.GetBool("Casbinv2.Enable") {
		if err := convertor.Map2Struct(variables.YmlAppConfig.Get("Casbinv2").(map[string]interface{}), &casbinCfg); err != nil {
			log.Fatalf("unable to init casbin v2 security with wrong yaml configuration")
		}
		InternalCasbin.NewCasbinContext(casbinCfg, variables.DBClient)
		variables.SecurityEnforcer = InternalCasbin.GetEnforcer()
	}

	// TODO initialize the internal cache container
	InternalContainer.New(variables.Logger)

	// TODO initialize and load the fasthttp websocket hub for fiber websocket
	if variables.YmlAppConfig.GetBool("Websocket.HubEnable") {
		variables.WebsocketHub = fiberws.GetHubInstance()
		variables.WebsocketHub.Run()
	}
}

func checkRequiredFolders() {
	if _, err := os.Stat(variables.BasePath + "/conf/application.yml"); err != nil {
		log.Fatalf("%s, details: %s", consts.ERRORS_CONFIG_APP_YML_FILE_NOT_EXISTS, err.Error())
	}
}

func migrateRequiredTables() {
	m := variables.DBClient.Migrator()

	var err error
	if err = m.AutoMigrate(model.Author{}); err != nil {
		goto errorHandle
	}

	if err = m.AutoMigrate(model.Blog{}); err != nil {
		goto errorHandle
	}

	if err = m.AutoMigrate(model.BlogCategory{}); err != nil {
		goto errorHandle
	}

	if err = m.AutoMigrate(model.BlogMetadata{}); err != nil {
		goto errorHandle
	}

	if err = m.AutoMigrate(model.BlogView{}); err != nil {
		goto errorHandle
	}

	if err = m.AutoMigrate(model.BlogClick{}); err != nil {
		goto errorHandle
	}

	if err = m.AutoMigrate(model.BlogVote{}); err != nil {
		goto errorHandle
	}
errorHandle:
	if err != nil {
		log.Fatalf("unable to auto migrate the required tables, error: [%s]\n", err.Error())
	}
}
