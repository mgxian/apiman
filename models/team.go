package models

import (
	//"errors"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type Team struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Name        string    `json:"name" gorm:"not null;unique"`
	CreatorID   uint      `json:"creator" gorm:"not null"`
	Description string    `json:"description"`
	AvatarUrl   string    `json:"avatar_url"`
	//DeletedAt   *time.Time `json:"-"`
	//Maintainers uint   `json:"maintainers" gorm:"not null"`
}

func (t *Team) AfterSave() (err error) {
	d := new(TeamIndex)
	copier.Copy(d, t)
	d.SearchType = "team"
	d.ID = strconv.Itoa(int(t.ID))
	err = BleveIndex.Index("user:"+d.ID, d)
	return
}

func (t *Team) AfterDelete() (err error) {
	err = BleveIndex.Delete("team:" + strconv.Itoa(int(t.ID)))
	return
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

type Projects struct {
	Project
	Creator string `json:"creator"`
	Team    string `json:"team"`
}

func (Projects) TableName() string {
	return "projects"
}

func GetTeamProjets(team_id uint) ([]*Projects, error) {
	tps := make([]*Projects, 0)
	err := db.Where("team_id = ?", team_id).Find(&tps).Error
	if err != nil {
		return nil, err
	}

	for _, tp := range tps {
		u, _ := GetUserByID(tp.CreatorID)
		t, _ := GetTeamByID(tp.TeamID)
		tp.Creator = u.Name
		tp.Team = t.Name
	}

	return tps, nil
}
