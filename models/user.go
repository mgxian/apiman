package models

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/pkg/jwt"
)

type User struct {
	//gorm.Model
	ID        uint       `json:"id" gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
	Name      string     `json:"name" gorm:"not null;unique" validate:"required"`
	Nickname  string     `json:"nickname" gorm:"not null" validate:"required"`
	Password  string     `json:"-" gorm:"not null"`
	AvatarUrl string     `json:"avatar_url"`
}

func CreateUser(u *User) error {
	u.Password, _ = getPassord(u.Password)
	err := db.Create(u).Error
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func GetUserByID(id uint) (*User, error) {
	u := new(User)
	err := db.First(u, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"id": id,
		}).Info("id not find in users")
		return nil, errors.New("id not find in users")
	}
	return u, nil
}

func GetUserByName(name string) (*User, error) {
	u := new(User)
	err := db.Where("name = ?", name).First(u).Error
	if err != nil {
		log.WithFields(log.Fields{
			"name": name,
		}).Info("name not find in users")
		return nil, errors.New("name not find in users")
	}
	return u, nil
}

func UpdateUser(u *User) error {
	//err := db.Save(u).Error
	//db.Model(&user).Updates(map[string]interface{}{"name": "hello", "age": 18, "actived": false})
	u.Password = ""
	err := db.Model(u).Updates(u).Error
	return err
}

func DeleteUserByUsername(name string) error {
	u, err := GetUserByName(name)
	err = db.Delete(u).Error
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
		log.Error("get password error")
		return false
	}
	u := new(User)
	err = db.Select("password").First(u, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"id": id,
		}).Info("no such user")
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
		log.Error("get password error")
		return errors.New("get password error")
	}
	u := new(User)
	u.Password = hash
	u.ID = id
	err = db.Model(u).Updates(u).Error
	if err != nil {
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
	if checkPassword(u.ID, password) {
		t, err := jwt.GetToken(u.Name, false)
		if err != nil {
			return "", err
		}
		data, err := jwt.ParseToken(t)
		log.WithFields(log.Fields{
			"token":     t,
			"username":  data.Name,
			"is_admin":  data.Admin,
			"expire_at": data.StandardClaims.ExpiresAt,
		}).Info("jwt parse")
		return t, nil
	}
	return "", errors.New("username or password is not right")
}
