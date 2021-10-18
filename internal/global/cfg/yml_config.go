package cfg

import (
	"github.com/fsnotify/fsnotify"
	"go.uber.org/zap"
	"log"
	"time"

	"github.com/LittleBenx86/Benlog/internal/global/consts"
	"github.com/LittleBenx86/Benlog/internal/repository/container"
	"github.com/LittleBenx86/Benlog/internal/utils/files/intf"

	"github.com/spf13/viper"
)

var (
	lastUpdateTime time.Time
)

func init() {
	lastUpdateTime = time.Now()
}

type ymlConfiguration struct {
	viper *viper.Viper
}

func (y *ymlConfiguration) isCached(key string) bool {
	if _, exists := container.GetInstance().IsExists(consts.PREFIX_CONFIG_KEY + key); exists {
		return true
	}
	return false
}

func (y *ymlConfiguration) cache(key string, value interface{}) bool {
	return container.GetInstance().Set(consts.PREFIX_CONFIG_KEY+key, value)
}

func (y *ymlConfiguration) getConfigValue(key string) interface{} {
	return container.GetInstance().Get(consts.PREFIX_CONFIG_KEY + key)
}

func (y *ymlConfiguration) clearCache() {
	container.GetInstance().FuzzyRemoveByPrefix(consts.PREFIX_CONFIG_KEY)
}

func (y *ymlConfiguration) ConfigFileUpdateListen() {
	/*
	   The viper package contains a bug:
	     When viper listens a file update event, it will trigger the related events for twice.
	     Bugid: https://github.com/spf13/viper/issues?q=OnConfigChange
	   Application level approach to avoid this question:
	     We set an inner global variable to record the file update timestamp. If these 2 callback functions' differ of timestamp
	   less than 1 second, we will treat it as the callback event instead of manually updating event.
	*/
	y.viper.OnConfigChange(func(changeEvent fsnotify.Event) {
		if time.Now().Sub(lastUpdateTime).Seconds() >= 1 {
			if changeEvent.Op.String() == "WRITE" {
				y.clearCache()
				lastUpdateTime = time.Now()
			}
		}
	})
	y.viper.WatchConfig()
}

// Clone
// Deep copy
func (y *ymlConfiguration) Clone(fileName string) intf.YmlConfig {
	var ymlConfigRef = *y
	var ymlConfigViperRef = *(y.viper)
	(&ymlConfigRef).viper = &ymlConfigViperRef
	(&ymlConfigRef).viper.SetConfigName(fileName)

	if err := (&ymlConfigRef).viper.ReadInConfig(); err != nil {
		log.Fatalf("%s, err: %+v\n", consts.ERRORS_CONFIG_INIT_ERR, zap.Error(err))
	}
	return &ymlConfigRef
}

func (y *ymlConfiguration) Get(key string) interface{} {
	if y.isCached(key) {
		return y.getConfigValue(key)
	}

	v := y.viper.Get(key)
	y.cache(key, v)
	return v
}

func (y *ymlConfiguration) GetString(key string) string {
	return y.Get(key).(string)
}

func (y *ymlConfiguration) GetBool(key string) bool {
	return y.Get(key).(bool)
}

func (y *ymlConfiguration) GetInt(key string) int {
	return y.Get(key).(int)
}

func (y *ymlConfiguration) GetInt32(key string) int32 {
	return y.Get(key).(int32)
}

func (y *ymlConfiguration) GetInt64(key string) int64 {
	return y.Get(key).(int64)
}

func (y *ymlConfiguration) GetFloat64(key string) float64 {
	return y.Get(key).(float64)
}

func (y *ymlConfiguration) GetDuration(key string) time.Duration {
	if y.isCached(key) {
		return y.getConfigValue(key).(time.Duration)
	}
	value := y.viper.GetDuration(key)
	y.cache(key, value)
	return value
}

func (y *ymlConfiguration) GetStringSlice(key string) []string {
	return y.Get(key).([]string)
}

// NewYmlConfigFactory
// Using the variable-length arguments means we can pass empty argument to this function.
// If the number of arguments more than one, we have to keep the first argument as useful one.
func NewYmlConfigFactory(rootPath string, fileName ...string) intf.YmlConfig {
	ymlConfig := viper.New()
	ymlConfig.AddConfigPath(rootPath + "/conf") // the location path of configuration files
	if len(fileName) == 0 {                     // set default to read file named "application"
		ymlConfig.SetConfigName("application")
	} else {
		ymlConfig.SetConfigName(fileName[0])
	}
	ymlConfig.SetConfigType("yml") // The file extension must be '.yml'

	if err := ymlConfig.ReadInConfig(); err != nil {
		log.Fatalf("%s, details: %s", consts.ERRORS_CONFIG_INIT_ERR, err.Error())
	}

	return &ymlConfiguration{
		ymlConfig,
	}
}
