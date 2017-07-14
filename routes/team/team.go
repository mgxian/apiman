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
	AvatarUrl   string `json:"avatar_url"`
	Creator     string `json:"creator" validate:"required,max=20"`
}

func CreateTeam(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

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

	username := tf.Creator
	u, _ := models.GetUserByName(username)

	if u == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "no such user",
		})
	}

	if u.Name != tokenInfo.Name && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	t.Creator = u.ID
	if err := models.CreateTeam(t); err != nil {
		log.Error(err.Error())
		return c.NoContent(http.StatusInternalServerError)
	}
	//fmt.Printf("%v\n", tf)
	//fmt.Printf("%v\n", t)
	copier.Copy(tf, t)
	tf.Creator = u.Name
	fmt.Printf("%v\n", tf)

	// add creator as the team's maintainer
	models.AddOrUpdateMember(tf.Name, username, models.Maintainer)

	return c.JSON(http.StatusCreated, tf)
}

func GetTeamByName(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	name := c.Param("teamname")
	t, err := models.GetTeamByName(name)
	if err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusNotFound)
	}

	tf := new(TeamForm)
	copier.Copy(tf, t)

	fmt.Printf("%v", tf)

	if u, err := models.GetUserByID(t.Creator); err == nil {
		tf.Creator = u.Name
		return c.JSON(http.StatusOK, tf)
	}

	return c.NoContent(http.StatusInternalServerError)
}

func AddOrUpdateTeamMember(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	tmf := new(struct {
		UserName string `json:"username" validate:"required,max=20"`
		Role     string `json:"role" validate:"required"`
		// maintainer member reader
	})

	if err = c.Bind(tmf); err != nil {
		fmt.Println(err)
	}

	if c.Request().Method == "POST" {
		if err := c.Validate(tmf); err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "request data is not right",
			})
		}
	}

	teamname := c.Param("teamname")

	operator := tokenInfo.Name

	flag := models.IsTeamMaintainer(teamname, operator)

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	var role int
	switch tmf.Role {
	case "maintainer":
		role = models.Maintainer
	case "member":
		role = models.Member
	case "reader":
		role = models.Reader
	default:
		role = models.Reader
	}

	if c.Request().Method == "PUT" {
		username := c.Param("username")
		err = models.AddOrUpdateMember(teamname, username, role)
	} else {
		err = models.AddOrUpdateMember(teamname, tmf.UserName, role)
	}

	if err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, tmf)

}

func RemoveTeamMember(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	teamname := c.Param("teamname")
	username := c.Param("username")

	operator := tokenInfo.Name
	flag := models.IsTeamMaintainer(teamname, operator)

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	err = models.RemoveMember(teamname, username)
	if err != nil {
		fmt.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

func GetTeamMembers(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	teamname := c.Param("teamname")
	users, _ := models.GetTeamMembers(teamname)
	for _, u := range users {
		fmt.Printf("team member %v\n", u)
	}

	return c.JSON(http.StatusOK, users)

}

func GetUserTeams(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	username := c.Param("username")
	teams, _ := models.GetUserTeams(username)
	for _, t := range teams {
		fmt.Printf("user team %v\n", t)
	}

	return c.JSON(http.StatusOK, teams)

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

	teamname := c.Param("teamname")
	username := tokenInfo.Name

	flag := models.IsTeamMaintainer(teamname, username)

	//fmt.Printf("flag------------%v", flag)

	if !flag && !tokenInfo.Admin {
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

	t, _ := models.GetTeamByName(teamname)
	//u, _ := models.GetUserByID(t.Creator)
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

	teamname := c.Param("teamname")
	username := tokenInfo.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	// delete all team member
	if err = models.RemoveAllMember(teamname); err == nil {
		fmt.Println("delete all memeber success")
	}

	if err = models.DeleteTeamByName(teamname); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}
