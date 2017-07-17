package models

import (
	//"errors"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type ApiGroup struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description"`
	Creator     uint      `json:"creator" gorm:"default 0"`
	Project     uint      `json:"project" gorm:"default 0"`
	//DeletedAt   *time.Time `json:"-"`
}

func CreateApiGroup(apigroup *ApiGroup) error {
	err := db.Create(apigroup).Error
	if err != nil {
		log.WithFields(log.Fields{
			"mysql": err.Error(),
		}).Info("create api group error")
		return err
	}

	return nil
}

func GetApiGroupByID(id uint) (*ApiGroup, error) {
	apigroup := new(ApiGroup)
	err := db.First(apigroup, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"mysql": err.Error(),
		}).Info("get api group error")
		return nil, err
	}

	return apigroup, nil
}

func UpdateApiGroup(apigroup *ApiGroup) error {
	err := db.Model(apigroup).Updates(apigroup).Error
	if err != nil {
		log.WithFields(log.Fields{
			"mysql": err.Error(),
		}).Info("update api group error")
		return err
	}

	return nil
}

func DeleteApiGroupByID(id uint) error {
	err := db.Where("id = ?", id).Delete(ApiGroup{}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"mysql": err.Error(),
		}).Info("delete api group error")
		return err
	}

	return nil
}
