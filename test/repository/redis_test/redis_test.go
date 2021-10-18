package redis_test

import (
	YmlFiles "github.com/LittleBenx86/Benlog/internal/global/cfg"
	"github.com/LittleBenx86/Benlog/internal/utils/convertor"
	"testing"

	InternalRedis "github.com/LittleBenx86/Benlog/internal/repository/redis"
	"github.com/mitchellh/mapstructure"
)

func Test_GetRedisConf(t *testing.T) {
	root := "/Users/benzheng/projects/Go/Benlog"
	ymlAppConfig := YmlFiles.NewYmlConfigFactory(root)
	redisConfMap := ymlAppConfig.Get("Redis").(map[string]interface{}) // map[string]interface{}, the field names are lower case
	t.Logf("%+v\n", redisConfMap)
	for k, v := range redisConfMap {
		t.Logf("%s : \"%+v\"", k, v)
	}

	var redisConf InternalRedis.PoolConfiguration
	err := mapstructure.Decode(redisConfMap, &redisConf)
	if err != nil {
		t.Fatal("map structure decode failed")
	}

	t.Logf("%+v\n", redisConf)
}

func Test_GetRedisClient(t *testing.T) {
	root := "/Users/benzheng/projects/Go/Benlog"
	tmpInitYml := YmlFiles.NewYmlConfigFactory(root)
	ymlAppConfig := tmpInitYml.Clone("application-" + tmpInitYml.GetString("Configuration.Active.App"))

	var cfg InternalRedis.PoolConfiguration
	_ = convertor.Map2Struct(ymlAppConfig.Get("Redis").(map[string]interface{}), &cfg)
	InternalRedis.NewRedisContext(cfg)
	instance := InternalRedis.GetInstance()
	t.Logf("%+v\n", instance)
}
