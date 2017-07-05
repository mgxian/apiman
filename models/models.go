package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/will835559313/apiman/pkg/setting"
)

var (
	DbCfg struct {
		Type, Host, Port, User, Password, Name string
	}
)

func loadConfigs() {
	//get database config
	setting.NewConfig()
	println("new cnfig finish")
	names := setting.Cfg.SectionStrings()
	println(names)
	sec := setting.Cfg.Section("database")
	println("get sec")
	dbtype := sec.Key("type").String()
	if dbtype == "mysql" {
		DbCfg.Host = sec.Key("host").String()
		DbCfg.Port = sec.Key("port").String()
		DbCfg.User = sec.Key("user").String()
		DbCfg.Password = sec.Key("password").String()
		DbCfg.Name = sec.Key("name").String()
		DbCfg.Type = "mysql"
	}
}

func GetDbConnection() (*gorm.DB, error) {
	loadConfigs()
	println("get config")
	connStr := DbCfg.User + ":" + DbCfg.Password + "@" + DbCfg.Host + "/" + DbCfg.Port + "/" + DbCfg.Name
	db, err := gorm.Open(DbCfg.Type, connStr)
	println(connStr)
	defer db.Close()
	return db, err
}
