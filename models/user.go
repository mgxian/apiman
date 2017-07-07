package models

import (
	"errors"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type User struct {
	//gorm.Model
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name" gorm:"not null;unique" validate:"required"`
	Nickname  string     `json:"nickname" gorm:"not null" validate:"required"`
	Password  string     `json:"-" gorm:"not null"`
	AvatarUrl string     `json:"avatar_url"`
}

func CreateUser(u *User) error {
	err := db.Create(u).Error
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func GetUserByID(id uint) (*User, error) {
	u := new(User)
	err := db.First(u, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"id": id,
		}).Info("id not find in users")
		return nil, errors.New("id not find in users")
	}
	return u, nil
}

func GetUserByName(name string) (*User, error) {
	u := new(User)
	err := db.Where("name = ?", name).First(u).Error
	if err != nil {
		log.WithFields(log.Fields{
			"name": name,
		}).Info("name not find in users")
		return nil, errors.New("name not find in users")
	}
	return u, nil
}

func UpdateUser(u *User) error {
	err := db.Save(u).Error
	return err
}
