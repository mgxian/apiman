package models

import (
	//"errors"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type Team struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Name        string    `json:"name" gorm:"not null;unique"`
	Creator     uint      `json:"creator" gorm:"not null"`
	Description string    `json:"description"`
	AvatarUrl   string    `json:"avatar_url"`
	//DeletedAt   *time.Time `json:"-"`
	//Maintainers uint   `json:"maintainers" gorm:"not null"`
}

func CreateTeam(t *Team) error {
	err := db.Create(t).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":   err.Error(),
			"team": *t,
		}).Error("create team error")
		return err
	}

	log.WithFields(log.Fields{
		"team": *t,
	}).Info("create team success")

	return nil
}

func GetTeamByName(name string) (*Team, error) {
	t := new(Team)
	err := db.Where("name = ?", name).First(t).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":   err.Error(),
			"name": name,
		}).Error("get team error")
		return nil, err
	}

	return t, nil
}

func GetTeamByID(id uint) (*Team, error) {
	t := new(Team)
	err := db.First(t, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("get team error")
		return nil, err
	}

	return t, nil
}

func UpdateTeam(t *Team) error {
	err := db.Model(t).Updates(t).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
		}).Error("update team error")
		return err
	}

	log.WithFields(log.Fields{
		"team": *t,
	}).Info("update team success")

	return nil
}

func DeleteTeamByName(name string) error {
	err := db.Where("name = ?", name).Delete(Team{}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":   err.Error(),
			"name": name,
		}).Error("delete team error")
		return err
	}

	return nil
}

func GetTeamProjets(team_id uint) ([]*Project, error) {
	tps := make([]*Project, 0)
	err := db.Where("team = ?", team_id).Find(&tps).Error
	if err != nil {
		return nil, err
	}

	return tps, nil
}
