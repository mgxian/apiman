package user

import (
	"fmt"
	//"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/models"
	//"gopkg.in/go-playground/validator.v9"
)

func CreateUser(c echo.Context) (err error) {
	u := new(models.User)
	if err = c.Bind(u); err != nil {
		fmt.Println(err)
		log.Info("user bind error")
		return
	}

	u, err = models.GetUserByName(u.Name)

	if u != nil {
		return c.NoContent(http.StatusConflict)
	}

	fmt.Println("build finish")
	if err = c.Validate(u); err != nil {
		return
	}

	fmt.Println("validator finish")
	if err = models.CreateUser(u); err != nil {
		log.Info("user create error")
		return c.NoContent(http.StatusInternalServerError)
	}
	log.WithFields(log.Fields{
		"username": u.Name,
	}).Info("user create username=" + u.Name)
	return c.JSONPretty(http.StatusCreated, u, "  ")
}

func GetUserByID(c echo.Context) (err error) {
	idstr := c.Param("id")
	idint, err := strconv.Atoi(idstr)
	u, err := models.GetUserByID(uint(idint))
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	if err = c.Validate(u); err != nil {
		log.Info("validator error")
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSONPretty(http.StatusOK, u, "  ")
}

func GetUserByName(c echo.Context) (err error) {
	name := c.Param("username")
	u := new(models.User)
	u, err = models.GetUserByName(name)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err = c.Validate(u); err != nil {
		log.Info("validator error")
		fmt.Println(err)
		//return c.NoContent(http.StatusNotFound)
		return err
	}
	return c.JSONPretty(http.StatusOK, u, "  ")
}

func UpdateUserByName(c echo.Context) error {
	name := c.Param("username")
	oldUser, err := models.GetUserByName(name)
	if err = c.Validate(oldUser); err != nil {
		log.Info("validator error")
		return c.NoContent(http.StatusNotFound)
	}
	newUser := new(models.User)
	c.Bind(newUser)
	newUser.ID = oldUser.ID
	err = models.UpdateUser(newUser)

	u, err := models.GetUserByName(name)
	//	if err = c.Validate(newUser); err != nil {
	//		log.Info("validator error")
	//		return c.NoContent(http.StatusInternalServerError)
	//	}
	if err != nil {
		log.Info(err)
	}
	return c.JSONPretty(http.StatusOK, u, "  ")
}

func DeleteUserByName(c echo.Context) error {
	name := c.Param("username")
	err := models.DeleteUserByUsername(name)
	if err != nil {
		log.Error("no such user")
	}
	return c.NoContent(http.StatusNoContent)
}

type Password struct {
	OldPassword string `json:"old_password" form:"old_password" query:"old_password"`
	NewPassword string `json:"new_password" form:"new_password" query:"new_password"`
}

func RestUserPassword(c echo.Context) error {
	name := c.Param("username")
	password := new(Password)
	if err := c.Bind(password); err != nil {
		fmt.Println(err)
	}
	u, err := models.GetUserByName(name)
	if u == nil {
		return c.NoContent(http.StatusNotFound)
	}
	err = models.RestPassword(password.OldPassword, password.NewPassword, u.ID)
	if err != nil {
		//data := map[string]string{}
		return c.JSONPretty(http.StatusBadRequest, echo.Map{
			"message": string(err.Error()),
		}, "  ")
	}
	return c.NoContent(http.StatusOK)
}

type Login struct {
	Name     string `json:name`
	Password string `json:password`
}

func GetToken(c echo.Context) error {
	login := new(Login)
	if err := c.Bind(login); err != nil {
		log.WithFields(log.Fields{
			"username": login.Name,
		}).Info("login fail")
		return c.JSONPretty(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		}, "  ")
	}
	if token, err := models.GetToken(login.Name, login.Password); err != nil {
		return c.JSONPretty(http.StatusBadRequest, echo.Map{
			"message": "用户名或密码错误",
		}, "  ")
	} else {
		return c.JSONPretty(http.StatusOK, echo.Map{
			"token": token,
		}, "  ")
	}
}
