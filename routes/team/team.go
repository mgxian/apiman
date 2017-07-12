package team

import (
	"fmt"
	//"io/ioutil"
	"net/http"
	//"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/models"
	//"github.com/will835559313/apiman/pkg/jwt"
	//"gopkg.in/go-playground/validator.v9"
)

type TeamForm struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Maintainers string `json:"maintainers"`
	AvatarUrl   string `json:"avatar_url"`
}

func CreateTeam(c echo.Context) error {
	t := new(models.Team)
	tb := new(TeamForm)
	if err := c.Bind(tb); err != nil {
		fmt.Println(err)
	}
	copier.Copy(t, tb)
	u, _ := models.GetUserByName(tb.Maintainers)
	if u == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "no such user",
		})
	}
	t.Maintainers = u.ID
	fmt.Printf("%v", tb)
	if err := models.CreateTeam(t); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}
	copier.Copy(tb, t)
	tb.ID = t.ID
	tb.Maintainers = u.Name
	return c.JSON(http.StatusOK, tb)
}
