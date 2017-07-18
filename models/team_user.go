package models

import (
	"errors"
	//"fmt"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

type TeamUser struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	TeamID    uint      `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	Role      uint      `gorm:"not null"`
}

const (
	Maintainer = 1
	Member     = 2
	Reader     = 3
)

func AddOrUpdateMember(teamname string, username string, role int) error {
	tu := new(TeamUser)
	t, _ := GetTeamByName(teamname)
	u, _ := GetUserByName(username)

	if t == nil || u == nil {
		log.WithFields(log.Fields{
			"team": teamname,
			"user": username,
		}).Error(err)
		return errors.New("get team or user error")
	}

	// select
	err := db.Where("team_id = ? and user_id = ?", t.ID, u.ID).First(tu).Error
	tu.UserID = u.ID
	tu.TeamID = t.ID
	tu.Role = uint(role)

	err = db.Save(tu).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"tu": *tu,
		}).Error("add or update member error")
	}

	return err
}

func RemoveMember(teamname, username string) error {
	t, _ := GetTeamByName(teamname)
	u, _ := GetUserByName(username)

	if t == nil || u == nil {
		log.WithFields(log.Fields{
			"team": teamname,
			"user": username,
		}).Error(err)
		return errors.New("get team or user error")
	}

	err := db.Where("team_id = ? and user_id = ?", t.ID, u.ID).Delete(TeamUser{}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":   err.Error(),
			"team": teamname,
			"user": username,
		}).Error("remove member error")
	}

	return err
}

func RemoveAllMember(teamname string) error {
	t, _ := GetTeamByName(teamname)
	if t == nil {
		log.WithFields(log.Fields{
			"team": teamname,
		}).Error(err)
		return errors.New("get team error")
	}

	//fmt.Println(t.Name)
	err := db.Where("team_id = ?", t.ID).Delete(TeamUser{}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":   err.Error(),
			"team": teamname,
		}).Error("remove all member error")
	}

	//fmt.Println(err)
	return err
}

type TeamMemberInfo struct {
	User
	Role string `json:"role"`
}

func GetTeamMembers(teamname string) ([]*TeamMemberInfo, error) {
	users := make([]*TeamMemberInfo, 0)
	tus := make([]*TeamUser, 0)

	t, _ := GetTeamByName(teamname)
	if t == nil {
		log.WithFields(log.Fields{
			"team": teamname,
		}).Error(err)
		return nil, errors.New("get team error")
	}

	err := db.Where("team_id = ?", t.ID).Find(&tus).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":   err.Error(),
			"team": teamname,
		}).Error("get team member error")
		return nil, errors.New("get team member error")
	}

	role := "reader"
	for _, tu := range tus {
		u, _ := GetUserByID(tu.UserID)
		switch tu.Role {
		case Maintainer:
			role = "maintainer"
		case Member:
			role = "member"
		case Reader:
			role = "reader"
		default:
		}
		users = append(users, &TeamMemberInfo{User: *u, Role: role})
	}
	return users, err
}

type UserTeams struct {
	Team
	Role string `json:"role"`
}

func GetUserTeams(username string) ([]*UserTeams, error) {
	userteams := make([]*UserTeams, 0)
	tus := make([]*TeamUser, 0)
	u, err := GetUserByName(username)
	if u == nil {
		log.WithFields(log.Fields{
			"user": username,
		}).Error(err)
		return nil, err
	}

	err = db.Where("user_id = ?", u.ID).Find(&tus).Error

	if err != nil {
		log.WithFields(log.Fields{
			"db":   err.Error(),
			"user": username,
		}).Error("get user team error")
		return nil, errors.New("get user team error")
	}

	role := "reader"
	for _, tu := range tus {
		t, _ := GetTeamByID(tu.TeamID)
		switch tu.Role {
		case Maintainer:
			role = "maintainer"
		case Member:
			role = "member"
		case Reader:
			role = "reader"
		default:
		}
		userteams = append(userteams, &UserTeams{Team: *t, Role: role})
	}

	return userteams, err
}

func IsTeamMaintainer(teamname, username string) bool {
	tu := new(TeamUser)
	t, _ := GetTeamByName(teamname)
	u, _ := GetUserByName(username)
	if t == nil || u == nil {
		log.WithFields(log.Fields{
			"team": teamname,
			"user": username,
		}).Error("get user or team error")
		return false
	}

	err := db.Where("team_id = ? and user_id = ? and role = ?", t.ID, u.ID, uint(Maintainer)).First(tu).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"tu": *tu,
		}).Error("check user maintainer role error")
		return false
	}

	return true
}

func IsTeamMember(teamname, username string) bool {
	tu := new(TeamUser)
	t, _ := GetTeamByName(teamname)
	u, _ := GetUserByName(username)
	if t == nil || u == nil {
		log.WithFields(log.Fields{
			"team": teamname,
			"user": username,
		}).Error("get user or team error")
		return false
	}

	err := db.Where("team_id = ? and user_id = ? and role = ?", t.ID, u.ID, uint(Member)).First(tu).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"tu": *tu,
		}).Error("check user member role error")
		return false
	}

	return true
}

func IsTeamReader(teamname, username string) bool {
	tu := new(TeamUser)
	t, _ := GetTeamByName(teamname)
	u, _ := GetUserByName(username)
	if t == nil || u == nil {
		log.WithFields(log.Fields{
			"team": teamname,
			"user": username,
		}).Error("get user or team error")
		return false
	}

	err := db.Where("team_id = ? and user_id = ? and role = ?", t.ID, u.ID, uint(Reader)).First(tu).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"tu": *tu,
		}).Error("check user reader role error")
		return false
	}

	return true
}
