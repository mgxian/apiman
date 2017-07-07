package user

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/models"
)

func CreateUser(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		fmt.Println(err)
		log.Info("user bind error")
	}
	if err := models.CreateUser(u); err != nil {
		log.Info("user create error")
		return c.NoContent(http.StatusInternalServerError)
	}
	log.WithFields(log.Fields{
		"username": u.Name,
	}).Info("user create username=" + u.Name)
	return c.JSONPretty(http.StatusCreated, u, "  ")
}
