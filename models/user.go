package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/will835559313/apiman/pkg/log"
)

type User struct {
	gorm.Model
	//ID       int64
	Name     string `gorm:"not null;unique"`
	Nickname string `gorm:"not null"`
	Password string `gorm:"not null"`
	Avatar   string
}

func (u *User) GetMyName() string {
	db, err := GetDbConnection()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("get db success")
	}
	fmt.Println(db)
	fmt.Println("exec sql")
	mysql.Exec("insert into users(name, nickname, password) value('will', 'will', 'will');")
	log, err := log.GetLogger()
	if err != nil {
		//fmt.Println("log errror")
		panic("get log error")
	}
	log.Info("hello db")
	return u.Nickname
}
