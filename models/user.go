package models

import (
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/mysql"
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
		db.Exec("select user();")
	}
	return u.Nickname
}
