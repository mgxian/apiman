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
	"github.com/will835559313/apiman/pkg/jwt"
	//"gopkg.in/go-playground/validator.v9"
)

type TeamForm struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" validate:"required,max=20"`
	Description string `json:"description" validate:"required,max=100"`
	Maintainers string `json:"maintainers" validate:"required,max=20"`
	AvatarUrl   string `json:"avatar_url"`
}

func CreateTeam(c echo.Context) error {
	t := new(models.Team)
	tf := new(TeamForm)
	if err := c.Bind(tf); err != nil {
		fmt.Println(err)
	}

	if err := c.Validate(tf); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "request data is not right",
		})
	}

	tf.ID = 0
	copier.Copy(t, tf)
	u, _ := models.GetUserByName(tf.Maintainers)
	if u == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "no such user",
		})
	}
	t.Maintainers = u.ID
	if err := models.CreateTeam(t); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}
	//fmt.Printf("%v\n", tf)
	//fmt.Printf("%v\n", t)
	copier.Copy(tf, t)
	//fmt.Printf("%v\n", tf)
	tf.Maintainers = u.Name
	return c.JSON(http.StatusCreated, tf)
}

func GetTeamByName(c echo.Context) error {
	name := c.Param("teamname")
	t, err := models.GetTeamByName(name)
	if err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusNotFound)
	}

	tf := new(TeamForm)
	copier.Copy(tf, t)

	if u, err := models.GetUserByID(t.ID); err == nil {
		tf.Maintainers = u.Name
		return c.JSON(http.StatusOK, tf)
	}

	return c.NoContent(http.StatusInternalServerError)
}

func UpdateTeamByName(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	name := c.Param("teamname")
	t, _ := models.GetTeamByName(name)
	u, _ := models.GetUserByID(t.Maintainers)

	if u.Name != tokenInfo.Name && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	tf := new(struct {
		Description string `json:"description" validate:"max=100"`
		AvatarUrl   string `json:"avatar_url"`
	})

	c.Bind(tf)

	if err := c.Validate(tf); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "request data is not right",
		})
	}

	copier.Copy(t, tf)

	if err := models.UpdateTeam(t); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, t)
}

func DeleteTeamByName(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	name := c.Param("teamname")
	t, _ := models.GetTeamByName(name)
	u, _ := models.GetUserByID(t.Maintainers)

	if u.Name != tokenInfo.Name && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	if err = models.DeleteTeamByName(name); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
