package security

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	MySQL "github.com/LittleBenx86/Benlog/internal/repository/mysql"
	"github.com/LittleBenx86/Benlog/internal/utils/files/intf"

	"github.com/casbin/casbin"
	"github.com/casbin/casbin/v2/model"
	GormAdapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var (
	dbType        string
	casbinContext *intf.YmlConfig
	initLock      sync.Mutex
	enforcer      *casbin.SyncedEnforcer
)

func Init(databaseType string, conf *intf.YmlConfig) {
	dbType = databaseType
	casbinContext = conf
}

func GetEnforcer() *casbin.SyncedEnforcer {
	if enforcer == nil {
		initLock.Lock()
		defer initLock.Unlock()

		e, err := newCasbinEnforcer()
		if err != nil {
			log.Fatalf(fmt.Sprintf("casbin v2 sync enforcer creation failed by error: %s", err.Error()))
		}
		enforcer = e
	}
	return enforcer
}

func newCasbinEnforcer() (*casbin.SyncedEnforcer, error) {
	var dbConn *gorm.DB
	switch strings.ToLower(dbType) {
	case "mysql":
		dbConn = MySQL.GetInstance()
	default:
		return nil, errors.New("unsupported database type to create connection and generate casbin table")
	}

	tbPrefix := (*casbinContext).GetString("Casbinv2.TablePrefix")
	tbName := (*casbinContext).GetString("Casbinv2.TableName")

	adapter, err := GormAdapter.NewAdapterByDBUseTableName(dbConn, tbPrefix, tbName)
	if err != nil {
		return nil, err
	}

	modelConf := (*casbinContext).GetString("Casbinv2.ModelConfig")
	m, err := model.NewModelFromString(modelConf)
	if err != nil {
		return nil, err
	}

	var enforcer *casbin.SyncedEnforcer
	if enforcer, err = casbin.NewSyncedEnforcerSafe(m, adapter); err != nil {
		return nil, err
	}

	_ = enforcer.LoadPolicy()
	autoLoad := (*casbinContext).GetDuration("Casbinv2.AutoLoadSeconds")
	enforcer.StartAutoLoadPolicy(autoLoad * time.Second)
	return enforcer, nil
}
