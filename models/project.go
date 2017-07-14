package models

import (
	//"errors"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type Project struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Creator     uint      `json:"creator" gorm:"default 0"`
	Team        uint      `json:"team" gorm:"default 0"`
	AvatarUrl   string    `json:"avatar_url"`
	//DeletedAt   *time.Time `json:"-"`
}

func CreateProject(t *Project) error {
	err := db.Create(t).Error
	if err != nil {
		log.Info(err.Error())
		return err
	}

	return nil
}

func GetProjectByID(id uint) (*Project, error) {
	p := new(Project)
	err := db.First(p, id).Error
	if err != nil {
		log.Info(err.Error())
		return nil, err
	}

	return p, nil
}

func UpdateProject(t *Project) error {
	err := db.Model(t).Updates(t).Error
	if err != nil {
		log.Info(err.Error())
		return err
	}

	return nil
}

func DeleteProjectByID(name string) error {
	err := db.Where("name = ?", name).Delete(Project{}).Error
	if err != nil {
		log.Info(err.Error())
		return err
	}

	return nil
}
