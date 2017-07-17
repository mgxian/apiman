package apigroup

import (
	"fmt"
	//"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
	//log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/models"
	"github.com/will835559313/apiman/pkg/jwt"
	//"gopkg.in/go-playground/validator.v9"
)

type ApiGroupForm struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" validate:"required,max=20"`
	Description string `json:"description" validate:"max=100"`
	Creator     string `json:"creator"`
	Project     string `json:"project"`
}

func CreateApiGroup(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	intstr, _ := strconv.Atoi(id)
	username := tokenInfo.Name

	p, _ := models.GetProjectByID(uint(intstr))
	t, _ := models.GetTeamByID(p.Team)
	u, _ := models.GetUserByName(username)

	teamname := t.Name
	//projectname := p.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag {
		flag = models.IsTeamMember(teamname, username)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	ag := new(models.ApiGroup)
	agf := new(ApiGroupForm)

	c.Bind(agf)

	if err := c.Validate(agf); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "request data is not right",
		})
	}

	agf.ID = 0
	copier.Copy(ag, agf)
	ag.Creator = u.ID
	ag.Project = p.ID

	if err := models.CreateApiGroup(ag); err == nil {
		agf.ID = ag.ID
		agf.Creator = u.Name
		agf.Project = p.Name
		return c.JSON(http.StatusCreated, agf)
	}

	return c.NoContent(http.StatusInternalServerError)
}

func GetApiGroupByID(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	intstr, _ := strconv.Atoi(id)

	a, _ := models.GetApiGroupByID(uint(intstr))
	p, _ := models.GetProjectByID(a.Project)
	t, _ := models.GetTeamByID(p.Team)
	u, _ := models.GetUserByID(a.Creator)

	username := tokenInfo.Name
	teamname := t.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag {
		flag = models.IsTeamMember(teamname, username)
	}

	if !flag {
		flag = models.IsTeamReader(teamname, username)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	if ag, err := models.GetApiGroupByID(uint(intstr)); err == nil {
		agf := new(ApiGroupForm)
		copier.Copy(agf, ag)
		agf.Creator = u.Name
		agf.Project = p.Name
		return c.JSON(http.StatusOK, agf)
	}

	return c.NoContent(http.StatusInternalServerError)
}

func UpdateApiGroupByID(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	intstr, _ := strconv.Atoi(id)
	username := tokenInfo.Name

	ag, _ := models.GetApiGroupByID(uint(intstr))
	p, _ := models.GetProjectByID(ag.Project)
	t, _ := models.GetTeamByID(p.Team)
	u, _ := models.GetUserByID(ag.Creator)

	teamname := t.Name
	//projectname := p.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag {
		flag = models.IsTeamMember(teamname, username)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	aguf := new(struct {
		Name        string `json:"name" validate:"required,max=20"`
		Description string `json:"description" validate:"max=100"`
	})

	if err := c.Bind(aguf); err != nil {
		fmt.Println(err)
	}

	if err := c.Validate(aguf); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "request data is not right",
		})
	}

	copier.Copy(ag, aguf)

	fmt.Printf("ag:--------%v-------------", ag)

	if err := models.UpdateApiGroup(ag); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	agf := new(ApiGroupForm)
	copier.Copy(agf, ag)
	agf.Creator = u.Name
	agf.Project = p.Name
	return c.JSON(http.StatusOK, agf)
}

func DeleteApiGroupByID(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	intstr, _ := strconv.Atoi(id)

	ag, _ := models.GetApiGroupByID(uint(intstr))
	p, _ := models.GetProjectByID(ag.Project)
	t, _ := models.GetTeamByID(p.Team)

	username := tokenInfo.Name
	teamname := t.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag {
		flag = models.IsTeamMember(teamname, username)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	if err := models.DeleteApiGroupByID(uint(intstr)); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
