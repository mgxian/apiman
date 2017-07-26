package models

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	//"fmt"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/pkg/jwt"
)

type User struct {
	//gorm.Model
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	//DeletedAt *time.Time `json:"-"`
	Name      string `json:"name" gorm:"not null;unique"`
	Nickname  string `json:"nickname" gorm:"not null"`
	Password  string `json:"-" gorm:"not null"`
	AvatarUrl string `json:"avatar_url"`
}

func CreateUser(u *User) error {
	u.Password, _ = getPassord(u.Password)
	err := db.Create(u).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":   err.Error(),
			"user": *u,
		}).Error("create user error")
		return err
	}

	log.WithFields(log.Fields{
		"user": *u,
	}).Info("create user success")

	return nil
}

func GetUserByID(id uint) (*User, error) {
	u := new(User)
	err := db.First(u, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("get user error")
		return nil, err
	}
	return u, nil
}

func GetUserByName(name string) (*User, error) {
	u := new(User)
	//fmt.Println(name)

	err := db.Where("name = ?", name).First(u).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":   err.Error(),
			"name": name,
		}).Error("get user error")
		//fmt.Println(err)
		return nil, err
	}
	//fmt.Printf("user: %v", u)
	return u, nil
}

func UpdateUser(u *User) error {
	//err := db.Save(u).Error
	//db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
	u.Password = ""
	err := db.Model(u).Updates(u).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":   err.Error(),
			"user": *u,
		}).Error("update user error")
	}

	log.WithFields(log.Fields{
		"user": *u,
	}).Info("update user success")

	return err
}

func DeleteUserByUsername(name string) error {
	u, err := GetUserByName(name)
	if u == nil {
		//fmt.Println(u)
		//fmt.Println(err)
		log.WithFields(log.Fields{
			"db":       err.Error(),
			"username": name,
		}).Error("get user error")
		return err
	}

	err = db.Delete(u).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":   err.Error(),
			"user": *u,
		}).Error("delete user error")
	}
	return err
}

func getHash(str string) (string, error) {
	hash := sha256.New()
	hash.Write([]byte(str))
	hashData := hash.Sum(nil)
	hashStr := hex.EncodeToString(hashData)
	return hashStr, nil
}

func getPassord(str string) (string, error) {
	return getHash(str)
}

func checkPassword(id uint, str string) bool {
	hash, err := getPassord(str)
	if err != nil {
		log.WithFields(log.Fields{
			"password": str,
		}).Error("gen password error")
		return false
	}

	u := new(User)
	err = db.Select("password").First(u, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"id":       id,
			"password": str,
		}).Error("get user password error")
		return false
	}

	if u.Password == hash {
		return true
	}

	return false
}

func setPassword(id uint, str string) error {
	hash, err := getPassord(str)
	if err != nil {
		log.WithFields(log.Fields{
			"password": str,
		}).Error("gen password error")
		return err
	}

	u := new(User)
	u.Password = hash
	u.ID = id
	err = db.Model(u).Updates(u).Error
	if err != nil {
		log.WithFields(log.Fields{
			"id":       id,
			"password": str,
		}).Error("set user password error")
		return err
	}

	return nil
}

func ChangeUserPassword(oldPassword, newPassword string, id uint) error {
	if !checkPassword(id, oldPassword) {
		return errors.New("passord is not right")
	}
	err := setPassword(id, newPassword)
	return err
}

func GetToken(name, password string) (string, error) {
	u, _ := GetUserByName(name)
	if u == nil {
		log.WithFields(log.Fields{
			"name":     name,
			"password": password,
		}).Error("get user error")

		return "", errors.New("get user error")
	}

	if checkPassword(u.ID, password) {
		t, err := jwt.GetToken(u.Name, false)
		if err != nil {
			return "", err
		}

		data, err := jwt.ParseToken(t)
		if err != nil {
			log.WithFields(log.Fields{
				"name":     name,
				"password": password,
			}).Error(err)
			return "", errors.New("parse token error")
		}

		log.WithFields(log.Fields{
			"token":     t,
			"username":  data.Name,
			"is_admin":  data.Admin,
			"expire_at": data.StandardClaims.ExpiresAt,
		}).Info("jwt parse success")
		return t, nil
	}

	return "", errors.New("username or password is not right")
}

type Teams struct {
	Team
	TeamID  uint   `json:"-"`
	UserID  uint   `json:"-"`
	RoleID  uint   `json:"-"`
	Role    string `json:"role"`
	Creator string `json:"creator"`
}

func (Teams) TableName() string {
	return "team_users"
}

func GetUserTeams(username string) ([]*Teams, error) {
	teams := make([]*Teams, 0)
	u, _ := GetUserByName(username)
	err := db.Where("user_id = ?", u.ID).Find(&teams).Error

	if err != nil {
		return nil, err
	}

	role := "reader"
	for _, team := range teams {
		t, _ := GetTeamByID(team.TeamID)
		u, _ := GetUserByID(t.CreatorID)
		switch team.RoleID {
		case Maintainer:
			role = "maintainer"
		case Member:
			role = "member"
		case Reader:
			role = "reader"
		default:
		}
		team.Role = role
		team.Creator = u.Name
		team.Name = t.Name
	}

	return teams, nil
}

func GetUserProjects(u_id uint) ([]*Projects, error) {
	ps := make([]*Projects, 0)
	tus := make([]*TeamUser, 0)

	if err := db.Where("user_id = ?", u_id).Find(&tus).Error; err != nil {
		return nil, err
	}

	for _, tu := range tus {
		ps_t := make([]*Projects, 0)
		err := db.Where("team_id = ?", tu.TeamID).Find(&ps_t).Error
		if err != nil {
			return nil, err
		}

		for _, p_t := range ps_t {
			ps = append(ps, p_t)
		}
	}

	for _, p := range ps {
		u, _ := GetUserByID(p.CreatorID)
		t, _ := GetTeamByID(p.TeamID)
		p.Creator = u.Name
		p.Team = t.Name
	}

	return ps, nil
}
