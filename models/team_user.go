package models

import (
	//"errors"
	"fmt"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	//log "github.com/sirupsen/logrus"
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

	// select
	err := db.Where("team_id = ? and user_id = ?", t.ID, u.ID).First(tu).Error
	tu.UserID = u.ID
	tu.TeamID = t.ID
	tu.Role = uint(role)

	err = db.Save(tu).Error

	return err
}

func RemoveMember(teamname, username string) error {
	t, _ := GetTeamByName(teamname)
	u, _ := GetUserByName(username)
	err := db.Where("team_id = ? and user_id = ?", t.ID, u.ID).Delete(TeamUser{}).Error
	return err
}

func RemoveAllMember(teamname string) error {
	t, _ := GetTeamByName(teamname)
	fmt.Println(t.Name)
	err := db.Where("team_id = ?", t.ID).Delete(TeamUser{}).Error
	fmt.Println(err)
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
	err := db.Where("team_id = ?", t.ID).Find(&tus).Error
	//fmt.Printf("%v", tus)
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
	//fmt.Println("---------first in--------------")
	userteams := make([]*UserTeams, 0)
	tus := make([]*TeamUser, 0)
	//fmt.Println("---------second in--------------")
	u, err := GetUserByName(username)
	if u == nil {
		fmt.Println("---------get user error--------------")
		return nil, err
	}
	//fmt.Println("---------third in--------------")
	err = db.Where("user_id = ?", u.ID).Find(&tus).Error
	//fmt.Println(err)
	//fmt.Printf("-----tus\n%v------\n", tus)

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
	fmt.Printf("userteams: %v", userteams)
	return userteams, err
}

func IsTeamMaintainer(teamname, username string) bool {
	tu := new(TeamUser)
	t, _ := GetTeamByName(teamname)
	u, _ := GetUserByName(username)
	err := db.Where("team_id = ? and user_id = ? and role = ?", t.ID, u.ID, uint(Maintainer)).First(tu).Error
	if err == nil {
		return true
	}
	//fmt.Printf("tu is %v", tu)
	fmt.Println(err)
	return false
}

func IsTeamMember(teamname, username string) bool {
	tu := new(TeamUser)
	t, _ := GetTeamByName(teamname)
	u, _ := GetUserByName(username)
	err := db.Where("team_id = ? and user_id = ? and role = ?", t.ID, u.ID, uint(Member)).First(tu).Error
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}

func IsTeamReader(teamname, username string) bool {
	tu := new(TeamUser)
	t, _ := GetTeamByName(teamname)
	u, _ := GetUserByName(username)
	err := db.Where("team_id = ? and user_id = ? and role = ?", t.ID, u.ID, uint(Reader)).First(tu).Error
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}
