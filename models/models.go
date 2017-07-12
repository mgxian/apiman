package models

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/will835559313/apiman/pkg/setting"
)

var (
	DbCfg struct {
		Type, Host, Port, User, Password, Name string
	}
	DB  *gorm.DB
	db  *gorm.DB
	err error
)

func loadConfigs() {
	sec := setting.Cfg.Section("database")
	//println("get sec")
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
	if DB != nil {
		return DB, nil
	}
	return nil, errors.New("DB is nil")
}

func Dbinit() {
	loadConfigs()
	//apiman:apiman@tcp(192.168.12.212:3306)/apiman
	connStr := DbCfg.User + ":" + DbCfg.Password + "@tcp(" + DbCfg.Host + ":" + DbCfg.Port + ")/" + DbCfg.Name
	connPara := "?charset=utf8&parseTime=True&loc=Local"
	fmt.Println("connstr: " + connStr)
	db, err = gorm.Open(DbCfg.Type, connStr+connPara)
	//defer db.Close()
	if err != nil {
		fmt.Println("connect error")
		fmt.Println(err)
	}
	DB = db
	//DB.Exec("insert into apiman_user(name, nickname, password) value('will', 'will', 'will');")
}

func DbMigrate() {
	mysql := DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8")
	mysql.AutoMigrate(&User{}, &Team{})
}
