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
	CreatorID   uint      `json:"creator" gorm:"default 0"`
	TeamID      uint      `json:"team" gorm:"default 0"`
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

func MigrateProjectByID(p_id, team_id uint) error {
	p := new(Project)
	p.ID = p_id
	p.TeamID = team_id
	err := db.Model(p).Updates(p).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":      err.Error(),
			"project": *p,
		}).Error("migrate project error")
		return err
	}

	log.WithFields(log.Fields{
		"project": *p,
	}).Info("migrate project success")

	return nil
}

func GetProjectApis(p_id uint) ([]*Apis, error) {
	//fmt.Println(p_id)
	apis := make([]*Apis, 0)
	err := db.Where("project_id = ?", p_id).Find(&apis).Error
	if err != nil {
		return nil, err
	}

	for _, api := range apis {
		u, _ := GetUserByID(api.CreatorID)
		api.Creator = u.Name
	}

	return apis, nil
}

type ApiGroups struct {
	ApiGroup
	Creator string `json:"creator"`
}

func (ApiGroups) TableName() string {
	return "api_groups"
}

func GetProjectApiGroups(p_id uint) ([]*ApiGroups, error) {
	api_groups := make([]*ApiGroups, 0)
	err := db.Where("project_id = ?", p_id).Find(&api_groups).Error
	if err != nil {
		return nil, err
	}

	for _, api_group := range api_groups {
		u, _ := GetUserByID(api_group.CreatorID)
		api_group.Creator = u.Name
	}

	return api_groups, nil
}
