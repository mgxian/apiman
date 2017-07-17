package project

import (
	"fmt"
	//"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/models"
	"github.com/will835559313/apiman/pkg/jwt"
	//"gopkg.in/go-playground/validator.v9"
)

type ProjectForm struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" validate:"required,max=100"`
	Description string `json:"description" validate:"max=100"`
	Creator     string `json:"creator"`
	Team        string `json:"team"`
	AvatarUrl   string `json:"avatar_url"`
}

func CreateProject(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	teamname := c.Param("teamname")
	username := tokenInfo.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	p := new(models.Project)
	pf := new(ProjectForm)
	if err := c.Bind(pf); err != nil {
		fmt.Println(err)
	}

	if err := c.Validate(pf); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "request data is not right",
		})
	}

	pf.ID = 0
	copier.Copy(p, pf)

	pf.Creator = username
	pf.Team = teamname

	u, _ := models.GetUserByName(username)
	t, _ := models.GetTeamByName(teamname)

	if u == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "no such user",
		})
	}

	if t == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "no such team",
		})
	}

	p.Creator = u.ID
	p.Team = t.ID
	if err := models.CreateProject(p); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}
	//fmt.Printf("%v\n", tf)
	//fmt.Printf("%v\n", t)
	copier.Copy(pf, p)
	//fmt.Printf("%v\n", tf)
	pf.Creator = u.Name
	pf.Team = t.Name
	return c.JSON(http.StatusCreated, pf)
}

func GetProjectByID(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	idint, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "id must be int",
			})
	}
	p, err := models.GetProjectByID(uint(idint))
	if err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusNotFound)
	}

	pf := new(ProjectForm)
	copier.Copy(pf, p)

	fmt.Printf("%v", pf)

	if u, err := models.GetUserByID(p.Creator); err == nil {
		pf.Creator = u.Name
	}

	if t, err := models.GetTeamByID(p.Team); err == nil {
		pf.Team = t.Name
	}

	if pf.Creator != "" && pf.Team != "" {
		return c.JSON(http.StatusOK, pf)
	}

	return c.NoContent(http.StatusInternalServerError)
}

func UpdateProjectByID(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	idint, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "id must be int",
			})
	}

	p, err := models.GetProjectByID(uint(idint))

	t, _ := models.GetTeamByID(p.Team)
	teamname := t.Name
	username := tokenInfo.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	puf := new(struct {
		Name        string `json:"name" validate:"required,max=100"`
		Description string `json:"description" validate:"max=100"`
		AvatarUrl   string `json:"avatar_url"`
	})

	//pf := new(ProjectForm)

	c.Bind(puf)

	if err := c.Validate(puf); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "request data is not right",
		})
	}

	avatarUrlOld := p.AvatarUrl
	copier.Copy(p, puf)
	if puf.AvatarUrl == "" {
		p.AvatarUrl = avatarUrlOld
	}

	if err := models.UpdateProject(p); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, p)
}

func DeleteProjectByID(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	idint, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "id must be int",
			})
	}
	p, err := models.GetProjectByID(uint(idint))

	t, _ := models.GetTeamByID(p.Team)
	teamname := t.Name
	username := tokenInfo.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	if err = models.DeleteProjectByID(uint(idint)); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
