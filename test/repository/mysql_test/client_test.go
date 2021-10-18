package mysql_test

import (
	"github.com/LittleBenx86/Benlog/internal/global/cfg"
	InternalDB "github.com/LittleBenx86/Benlog/internal/repository/mysql"
	"github.com/LittleBenx86/Benlog/internal/utils/convertor"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	//_ "github.com/LittleBenx86/Benlog/internal/bootstrap"
)

func Test_GetMySQLClient(t *testing.T) {
	rootPath := "/Users/benzheng/Projects/Go/Benlog"
	tmpInitYml := cfg.NewYmlConfigFactory(rootPath)
	ymlAppConfig := tmpInitYml.Clone("application-" + tmpInitYml.GetString("Configuration.Active.App"))
	var cfg InternalDB.DBConfigParams
	t.Logf("%+v\n", ymlAppConfig.Get("Gormv2.MySQL"))
	_ = convertor.Map2Struct(ymlAppConfig.Get("Gormv2.MySQL").(map[string]interface{}), &cfg)
	InternalDB.NewDBContext(cfg, nil)
	t.Logf("%+v\n", cfg)
	t.Logf("%+v\n", InternalDB.GetInstance())
}

func Test_GormWithSqlMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	defer func() {
		err := db.Close()
		if err != nil {
			return
		}
	}()

	if err != nil {
		t.Fatalf("init sqlmock failed, err %v\n", err)
	}

	t.Logf("%+v\n", db)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE abc").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

}
