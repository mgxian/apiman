package models

import (
	//"errors"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type Team struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	//DeletedAt   *time.Time `json:"-"`
	Name        string `json:"name" gorm:"not null;unique"`
	Description string `json:"description"`
	Maintainers uint   `json:"maintainers" gorm:"not null"`
	AvatarUrl   string `json:"avatar_url"`
}

func CreateTeam(t *Team) error {
	err := db.Create(t).Error
	if err != nil {
		log.Info(err.Error())
		return err
	}

	return nil
}

func GetTeamByName(name string) (*Team, error) {
	t := new(Team)
	err := db.Where("name = ?", name).First(t).Error
	if err != nil {
		log.Info(err.Error())
		return nil, err
	}

	return t, nil
}

func UpdateTeam(t *Team) error {
	err := db.Model(t).Updates(t).Error
	if err != nil {
		log.Info(err.Error())
		return err
	}

	return nil
}

func DeleteTeamByName(name string) error {
	err := db.Where("name = ?", name).Delete(Team{}).Error
	if err != nil {
		log.Info(err.Error())
		return err
	}

	return nil
}
