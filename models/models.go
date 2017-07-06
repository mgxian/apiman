package models

import (
	"fmt"

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
	//setting.NewConfig()
	//println("new cnfig finish")
	//names := setting.Cfg.SectionStrings()
	//println(names)
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
	//apiman:apiman@tcp(192.168.12.212:3306)/apiman
	connStr := DbCfg.User + ":" + DbCfg.Password + "@tcp(" + DbCfg.Host + ":" + DbCfg.Port + ")/" + DbCfg.Name
	println("connstr: " + connStr)
	db, err := gorm.Open(DbCfg.Type, connStr)
	defer db.Close()
	if err != nil {
		println("connect error")
		fmt.Println(err)
		return db, err
	}
	return db, nil
}
