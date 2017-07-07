package models

import (
	//"fmt"
	"time"

	//"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type User struct {
	//gorm.Model
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name" gorm:"not null;unique"`
	Nickname  string     `json:"nickname" gorm:"not null"`
	Password  string     `json:"-" gorm:"not null"`
	AvatarUrl string     `json:"avatar_url"`
}

func CreateUser(u *User) error {
	err = db.Create(u).Error
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (u *User) GetMyName() string {
	log.Info("get my name")
	return u.Nickname
}
