package mysql

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/LittleBenx86/Benlog/internal/global/consts"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

type DBConfigDetail struct {
	Host               string        `json:"Host" mapstructure:"host"`
	Port               int           `json:"Port" mapstructure:"port"`
	DataBase           string        `json:"DataBase" mapstructure:"database"`
	TablePrefix        string        `json:"TablePrefix" mapstructure:"tableprefix"`
	User               string        `json:"User" mapstructure:"user"`
	Password           string        `json:"Password" mapstructure:"password"`
	Charset            string        `json:"Charset" mapstructure:"charset"`
	MaxIdleConns       int           `json:"SetMaxIdleConns" mapstructure:"setmaxidleconns"`
	MaxOpenConns       int           `json:"SetMaxOpenConns" mapstructure:"setmaxopenconns"`
	ConnMaxLiveTime    time.Duration `json:"SetConnMaxLiveTime" mapstructure:"setconnmaxlivetime"`
	ReconnInterval     int           `json:"ReconnectInterval" mapstructure:"reconnectinterval,omitempty"`
	PingFailRetryTimes int           `json:"PingFailRetryTimes" mapstructure:"pingfailretrytimes,omitempty"`
}

type DBConfigParams struct {
	InitDBStartup          bool           `json:"IsInitDBStartup" mapstructure:"isinitdbstartup"`
	ReadWriteSplitting     bool           `json:"IsRWSplittingEnable" mapstructure:"isrwsplittingenable"`
	GlobalSlowSQLThreshold time.Duration  `json:"SlowThreshold" mapstructure:"slowthreshold"`
	Write                  DBConfigDetail `json:"Write" mapstructure:"write"`
	Read                   DBConfigDetail `json:"Read" mapstructure:"read"`
}

func (d *DBConfigParams) GetDetailByRWType(rw string) (*DBConfigDetail, error) {
	rw = strings.ToLower(rw)
	if rw == "write" {
		return &d.Write, nil
	} else if rw == "read" {
		return &d.Read, nil
	}
	return nil, errors.New("unknown read or write mode")
}

type DBContext struct {
	Cfg    DBConfigParams
	Logger *zap.Logger
}

var (
	lock      sync.Mutex
	instance  *gorm.DB
	dbContext *DBContext
)

const (
	SUPPORTED_DB_DRIVER_MYSQL = "mysql"
)

func NewDBContext(cfg DBConfigParams, l *zap.Logger) {
	dbContext = &DBContext{
		Cfg:    cfg,
		Logger: l,
	}
}

func GetInstance() *gorm.DB {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()

		if instance == nil {
			d, err := getSqlDriver("MySQL", dbContext.Cfg.ReadWriteSplitting)
			if err != nil {
				log.Fatalf(fmt.Sprintf("unable to generate mysql client by error: %s", err.Error()))
			}
			instance = d
		}
	}
	return instance
}

func getSqlDriver(dbType string, rwsMode bool) (*gorm.DB, error) {
	var dbDialect gorm.Dialector
	if val, err := getDBDialect(dbType, "Write"); err != nil { // write database must present first
		msg := fmt.Sprintf("%s details: %s", consts.ERRORS_DB_DIALECT_INIT_ERR, dbType)
		dbContext.Logger.Error(msg)
	} else {
		dbDialect = val
	}

	gormDb, err := gorm.Open(dbDialect, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logInterceptor(dbType),
	})
	if err != nil {
		// gorm sql driver init failed
		return nil, err
	}

	// If the read-write splitting is enabled, the readonly database's resource, read and replicas should be configured
	if rwsMode {
		if dial, err := getDBDialect(dbType, "Read"); err != nil {
			msg := fmt.Sprintf("%s, database type: %s", consts.ERRORS_DB_DIALECT_INIT_ERR, dbType)
			dbContext.Logger.Error(msg, zap.Error(err))
		} else {
			dbDialect = dial
		}

		resolverCfg := dbresolver.Config{
			Replicas: []gorm.Dialector{dbDialect}, // For query only operation
			Policy:   dbresolver.RandomPolicy{},   // satisfy the sources or replicas load balance policy
		}

		err = gormDb.Use(dbresolver.
			Register(resolverCfg).
			SetConnMaxIdleTime(30 * time.Second).
			SetConnMaxLifetime(dbContext.Cfg.Read.ConnMaxLiveTime * time.Second).
			SetMaxIdleConns(dbContext.Cfg.Read.MaxIdleConns).
			SetMaxOpenConns(dbContext.Cfg.Read.MaxOpenConns))

		if err != nil {
			return nil, err
		}
	}

	// If query result is empty, we mask the gorm v2 will occur bug
	// issue: https://github.com/go-gorm/gorm/issues/3789
	_ = gormDb.Callback().Query().Before("gorm:query").
		Register("disable_raise_record_not_found", func(d *gorm.DB) {
			d.Statement.RaiseErrorOnNotFound = false
		})

	// Set connections pool for main connection
	if rawDb, err := gormDb.DB(); err != nil {
		return nil, err
	} else {
		rawDb.SetConnMaxIdleTime(30 * time.Second)
		rawDb.SetConnMaxLifetime(dbContext.Cfg.Read.ConnMaxLiveTime * time.Second)
		rawDb.SetMaxIdleConns(dbContext.Cfg.Read.MaxIdleConns)
		rawDb.SetMaxOpenConns(dbContext.Cfg.Read.MaxOpenConns)
		return gormDb, nil
	}
}

func logInterceptor(dbType string) gormLogger.Interface {
	return NewCustomGormDBLog(
		dbType,
		SetInfoFormat("[info] %s %s \n"),
		SetWarnFormat("[warn] %s %s \n"),
		SetErrorFormat("[error] %s %s \n"),
		SetTraceFormat("[trace] %s %s [%.3fms] [rows:%v] %s \n"),
		SetTraceWarnFormat("[traceWarn] %s %s %s [%.3fms] [rows:%v] %s \n"),
		SetTraceErrorFormat("[traceError] %s %s %s [%.3fms] [rows:%v] %s \n"))
}

func getDBDialect(dbType string, rw string) (gorm.Dialector, error) {
	var dbDialect gorm.Dialector
	dsn := getDsn(dbType, rw)
	switch strings.ToLower(dbType) {
	case SUPPORTED_DB_DRIVER_MYSQL:
		dbDialect = mysql.New(mysql.Config{
			DSN:                       dsn,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		})
	default:
		return nil, errors.New(consts.ERRORS_DB_DRIVER_UNSUPPORTED)
	}
	return dbDialect, nil
}

func getDsn(dbType string, rw string) string {
	detail, err := dbContext.Cfg.GetDetailByRWType(rw)
	if err != nil {
		log.Fatalf("Unsupported read-write splitting type [%s]", rw)
	}

	switch strings.ToLower(dbType) {
	case SUPPORTED_DB_DRIVER_MYSQL:
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			detail.User, detail.Password, detail.Host, detail.Port, detail.DataBase, detail.Charset)
	}
	return ""
}
