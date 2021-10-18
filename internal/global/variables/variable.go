package variables

import (
	FileIntf "github.com/LittleBenx86/Benlog/internal/utils/files/intf"
	"github.com/LittleBenx86/Benlog/internal/utils/websocket/ginws"
	"github.com/casbin/casbin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// We have to check the global variables support the CAS features.

var (
	BasePath         string // The project root path
	YmlAppConfig     FileIntf.YmlConfig
	Logger           *zap.Logger
	DBInstance       *gorm.DB
	RedisInstance    *redis.Client
	SecurityInstance *casbin.SyncedEnforcer
	WebsocketHub     *ginws.Hub
)
