package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/will835559313/apiman/pkg/log"
)

type User struct {
	gorm.Model
	ID       int64
	Name     string
	Nickname string
	Password string
}

func (u *User) GetMyName() string {
	db, err := GetDbConnection()
	if err != nil {
		println("exec sql")
		fmt.Println(err)
	}
	db.Exec("select user();")
	if err != nil {
		fmt.Println(err)
	}
	log, err := log.GetLogger()
	if err != nil {
		//fmt.Println("log errror")
		panic("get log error")
	}
	log.Info("hello db")
	return u.Nickname
}
