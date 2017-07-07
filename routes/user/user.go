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
	if err = c.Validate(newUser); err != nil {
		log.Info("validator error")
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSONPretty(http.StatusOK, newUser, "  ")
}
