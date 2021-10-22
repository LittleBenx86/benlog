package security

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	GormAdapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

var (
	once      sync.Once
	enforcer  *casbin.SyncedEnforcer
	casbinCtx *CasbinContext
)

type CasbinContext struct {
	Cfg      CasbinConfiguration
	BbClient *gorm.DB
}

type CasbinConfiguration struct {
	Enable                 bool     `json:"enable" mapstructure:"enable"`
	FirstInit              bool     `json:"first_init" mapstructure:"firstinit"`
	AutoLoadPoliciesPeriod int      `json:"auto_load_policies_period" mapstructure:"autoloadpoliciesperiod"`
	TablePrefix            string   `json:"table_prefix" mapstructure:"tableprefix"`
	TableName              string   `json:"table_name" mapstructure:"tablename"`
	ModelCfgContent        string   `json:"model_cfg_content" mapstructure:"modelconfig"`
	Policies               []Policy `json:"policies" mapstructure:"policies"`
}

type Policy struct {
	Ptype string `json:"ptype" mapstructure:"ptype"`
	V0    string `json:"v0" mapstructure:"v0"`
	V1    string `json:"v1" mapstructure:"v1"`
	V2    string `json:"v2" mapstructure:"v2"`
	V3    string `json:"v3" mapstructure:"v3"`
	V4    string `json:"v4" mapstructure:"v4"`
	V5    string `json:"v5" mapstructure:"v5"`
}

func NewCasbinContext(cfg CasbinConfiguration, client *gorm.DB) {
	casbinCtx = &CasbinContext{
		Cfg:      cfg,
		BbClient: client,
	}
}

func GetCasbinTable() string {
	return casbinCtx.Cfg.TablePrefix + "_" + casbinCtx.Cfg.TableName
}

func GetEnforcer() *casbin.SyncedEnforcer {
	once.Do(func() {
		e, err := newCasbinEnforcer()
		if err != nil {
			log.Fatalf(fmt.Sprintf("casbin v2 sync enforcer creation failed by error: %s", err.Error()))
		}
		enforcer = e
	})

	return enforcer
}

func newCasbinEnforcer() (*casbin.SyncedEnforcer, error) {
	adapter, err := GormAdapter.NewAdapterByDBUseTableName(
		casbinCtx.BbClient, casbinCtx.Cfg.TablePrefix, casbinCtx.Cfg.TableName,
	) // It will create the casbin rules table automatically.
	if err != nil {
		return nil, err
	}

	if casbinCtx.Cfg.FirstInit && isPoliciesEmpty() {
		storeInitPolicies()
	}

	m, err := model.NewModelFromString(casbinCtx.Cfg.ModelCfgContent)
	if err != nil {
		return nil, err
	}

	var enforcer *casbin.SyncedEnforcer
	if enforcer, err = casbin.NewSyncedEnforcer(m, adapter); err != nil {
		return nil, err
	}

	_ = enforcer.LoadPolicy()
	autoLoad := time.Duration(casbinCtx.Cfg.AutoLoadPoliciesPeriod)
	enforcer.StartAutoLoadPolicy(autoLoad * time.Second)
	return enforcer, nil
}

func storeInitPolicies() {
	policies := make([]map[string]interface{}, 0, 10)
	for _, p := range casbinCtx.Cfg.Policies {
		policies = append(policies, map[string]interface{}{
			"ptype": p.Ptype,
			"v0":    p.V0,
			"v1":    p.V1,
			"v2":    p.V2,
			"v3":    p.V3,
			"v4":    p.V4,
			"v5":    p.V5,
		})
	}

	casbinCtx.BbClient.
		Table(casbinCtx.Cfg.TablePrefix+"_"+casbinCtx.Cfg.TableName).
		CreateInBatches(policies, 50)
}

func isPoliciesEmpty() bool {
	var count int64
	casbinCtx.BbClient.
		Table(casbinCtx.Cfg.TablePrefix + "_" + casbinCtx.Cfg.TableName).
		Count(&count)
	return count == 0
}
