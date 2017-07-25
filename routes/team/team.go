package team

import (
	//"fmt"
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

	//fmt.Println(tokenInfo.Name)

	t := new(models.Team)
	tf := new(TeamForm)
	if err := c.Bind(tf); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if err := c.Validate(tf); err != nil {
		//fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if t, _ := models.GetTeamByName(tf.Name); t != nil {
		return c.NoContent(http.StatusConflict)
	}

	tf.ID = 0
	copier.Copy(t, tf)

	username := tf.Creator
	u, _ := models.GetUserByName(username)
	if u == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "用户不存在",
		})
	}

	if u.Name != tokenInfo.Name && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	t.Creator = u.ID
	if err := models.CreateTeam(t); err != nil {
		log.WithFields(log.Fields{
			"team": *t,
		}).Error("create team error")
		return c.NoContent(http.StatusInternalServerError)
	}

	copier.Copy(tf, t)
	tf.Creator = u.Name
	//fmt.Printf("%v\n", tf)

	// add creator as the team's maintainer
	err = models.AddOrUpdateMember(tf.Name, username, models.Maintainer)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"team":     *t,
		"operator": tokenInfo.Name,
	}).Info("create team success")

	return c.JSON(http.StatusCreated, tf)
}

func GetTeamByName(c echo.Context) error {
	_, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	name := c.Param("teamname")
	t, err := models.GetTeamByName(name)
	if err != nil {
		//fmt.Println(err)
		return c.NoContent(http.StatusNotFound)
	}

	tf := new(TeamForm)
	copier.Copy(tf, t)

	//fmt.Printf("%v", tf)

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

	//fmt.Println(tokenInfo.Name)

	tmf := new(struct {
		UserName string `json:"username" validate:"required,max=20"`
		Role     string `json:"role" validate:"required"`
		// maintainer member reader
	})

	if err = c.Bind(tmf); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if c.Request().Method == "POST" {
		if err := c.Validate(tmf); err != nil {
			//fmt.Println(err)
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "请求数据错误",
			})
		}
	}

	teamname := c.Param("teamname")
	operator := tokenInfo.Name

	flag := models.IsTeamMaintainer(teamname, operator)

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
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
		log.WithFields(log.Fields{
			"team": teamname,
			"user": c.Param("username"),
		}).Error("add or update team member error")
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"team":     teamname,
		"user":     c.Param("username"),
		"operator": tokenInfo.Name,
	}).Info("add or update team memeber success")

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

	//fmt.Println(tokenInfo.Name)

	teamname := c.Param("teamname")
	username := c.Param("username")

	operator := tokenInfo.Name
	flag := models.IsTeamMaintainer(teamname, operator)

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	err = models.RemoveMember(teamname, username)
	if err != nil {
		log.WithFields(log.Fields{
			"team": teamname,
			"user": username,
		}).Error("remove team member error")
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"team":     teamname,
		"user":     username,
		"operator": tokenInfo.Name,
	}).Info("remove team memeber success")

	return c.NoContent(http.StatusNoContent)
}

func GetTeamMembers(c echo.Context) error {
	_, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	teamname := c.Param("teamname")
	users, _ := models.GetTeamMembers(teamname)

	if users == nil {
		log.WithFields(log.Fields{
			"team": teamname,
		}).Error("get team members error")
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, users)
}

func GetUserTeams(c echo.Context) error {
	_, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	username := c.Param("username")
	teams, _ := models.GetUserTeams(username)

	if teams == nil {
		log.WithFields(log.Fields{
			"user": username,
		}).Error("get user teams error")
		return c.NoContent(http.StatusInternalServerError)
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

	//fmt.Println(tokenInfo.Name)

	teamname := c.Param("teamname")
	username := tokenInfo.Name

	flag := models.IsTeamMaintainer(teamname, username)

	//fmt.Printf("flag------------%v", flag)

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	tf := new(struct {
		Description string `json:"description" validate:"max=100"`
		AvatarUrl   string `json:"avatar_url"`
	})

	if err := c.Bind(tf); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if err := c.Validate(tf); err != nil {
		//fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	t, _ := models.GetTeamByName(teamname)
	if t == nil {
		return c.NoContent(http.StatusNotFound)
	}

	if tf.AvatarUrl == "" {
		tf.AvatarUrl = t.AvatarUrl
	}
	copier.Copy(t, tf)

	if err := models.UpdateTeam(t); err != nil {
		log.WithFields(log.Fields{
			"team": teamname,
		}).Error("update team info error")
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"team":     *t,
		"operator": tokenInfo.Name,
	}).Info("update team info success")

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

	//fmt.Println(tokenInfo.Name)

	teamname := c.Param("teamname")
	username := tokenInfo.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	// delete all team member
	if err = models.RemoveAllMember(teamname); err != nil {
		log.WithFields(log.Fields{
			"team": teamname,
		}).Error("delete all team member error")
		return c.NoContent(http.StatusInternalServerError)
	}

	if err = models.DeleteTeamByName(teamname); err != nil {
		log.WithFields(log.Fields{
			"team": teamname,
		}).Error("delete team error")
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"team":     teamname,
		"operator": tokenInfo.Name,
	}).Info("delete team success")

	return c.NoContent(http.StatusNoContent)
}

func GetTeamMember(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	operator := tokenInfo.Name
	username := c.Param("username")
	teamname := c.Param("teamname")

	t, _ := models.GetTeamByName(teamname)
	if t == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "team不存在",
		})
	}

	flag := models.IsTeamMaintainer(teamname, operator)

	if !flag {
		flag = models.IsTeamMember(teamname, operator)
	}

	if !flag {
		flag = models.IsTeamReader(teamname, operator)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	tm, _ := models.GetTeamMemberByID(teamname, username)

	if tm == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "成员不存在",
		})
	}

	return c.JSON(http.StatusOK, tm)
}

func GetTeamProjets(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	operator := tokenInfo.Name
	teamname := c.Param("teamname")

	t, _ := models.GetTeamByName(teamname)
	if t == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "team不存在",
		})
	}

	flag := models.IsTeamMaintainer(teamname, operator)

	if !flag {
		flag = models.IsTeamMember(teamname, operator)
	}

	if !flag {
		flag = models.IsTeamReader(teamname, operator)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	tps, _ := models.GetTeamProjets(t.ID)
	if len(tps) == 0 {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, tps)
}
