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

func CreateProject(p *Project) error {
	err := db.Create(p).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":      err.Error(),
			"project": *p,
		}).Error("create project error")
		return err
	}

	log.WithFields(log.Fields{
		"project": *p,
	}).Info("create project success")

	return nil
}

func GetProjectByID(id uint) (*Project, error) {
	p := new(Project)
	err := db.First(p, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("get project error")
		return nil, err
	}

	return p, nil
}

func UpdateProject(p *Project) error {
	err := db.Model(p).Updates(p).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":      err.Error(),
			"project": *p,
		}).Error("update project error")
		return err
	}

	log.WithFields(log.Fields{
		"project": *p,
	}).Info("update project success")

	return nil
}

func DeleteProjectByID(id uint) error {
	err := db.Where("id = ?", id).Delete(Project{}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("delete project error")
		return err
	}

	return nil
}
