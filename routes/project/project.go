package project

import (
	//"fmt"
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
	Name        string `json:"name" validate:"required,max=20"`
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

	//fmt.Println(tokenInfo.Name)

	teamname := c.Param("teamname")
	username := tokenInfo.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusForbidden,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	p := new(models.Project)
	pf := new(ProjectForm)
	if err := c.Bind(pf); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if err := c.Validate(pf); err != nil {
		//fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
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
			"message": "用户不存在",
		})
	}

	if t == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Team不存在",
		})
	}

	p.CreatorID = u.ID
	p.TeamID = t.ID
	if err := models.CreateProject(p); err != nil {
		log.WithFields(log.Fields{
			"project":  *p,
			"operator": username,
		}).Error("create project fail")
		return c.NoContent(http.StatusInternalServerError)
	}

	copier.Copy(pf, p)
	pf.Creator = u.Name
	pf.Team = t.Name

	log.WithFields(log.Fields{
		"project":  *p,
		"operator": username,
	}).Info("create project success")

	return c.JSON(http.StatusCreated, pf)
}

func GetProjectByID(c echo.Context) error {
	_, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	idint, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "ID必须是整数",
			})
	}

	p, err := models.GetProjectByID(uint(idint))
	if err != nil {
		//fmt.Println(err)
		return c.NoContent(http.StatusNotFound)
	}

	pf := new(ProjectForm)
	copier.Copy(pf, p)

	//fmt.Printf("%v", pf)

	if u, err := models.GetUserByID(p.CreatorID); err == nil {
		pf.Creator = u.Name
	} else {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "用户不存在",
		})
	}

	if t, err := models.GetTeamByID(p.TeamID); err == nil {
		pf.Team = t.Name
	} else {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Team不存在",
		})
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

	//fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	idint, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "ID必须是整数",
			})
	}

	p, err := models.GetProjectByID(uint(idint))

	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "project不存在",
		})
	}

	t, _ := models.GetTeamByID(p.TeamID)
	if t == nil {
		log.WithFields(log.Fields{
			"project": *p,
		}).Error("project's team is not exist")

		return c.NoContent(http.StatusInternalServerError)
	}

	teamname := t.Name
	username := tokenInfo.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusForbidden,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	puf := new(struct {
		Name        string `json:"name" validate:"required,max=20"`
		Description string `json:"description" validate:"max=100"`
		AvatarUrl   string `json:"avatar_url"`
	})

	//pf := new(ProjectForm)

	if err := c.Bind(puf); err != nil {
		//fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if err := c.Validate(puf); err != nil {
		//fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	avatarUrlOld := p.AvatarUrl
	copier.Copy(p, puf)
	if puf.AvatarUrl == "" {
		p.AvatarUrl = avatarUrlOld
	}

	if err := models.UpdateProject(p); err != nil {
		log.WithFields(log.Fields{
			"project":  *p,
			"operator": tokenInfo.Name,
		}).Error("update project info fail")
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"project":  *p,
		"operator": tokenInfo.Name,
	}).Info("update project info success")

	u, _ := models.GetUserByID(p.CreatorID)
	pf := new(ProjectForm)
	copier.Copy(pf, p)
	pf.Team = teamname
	pf.Creator = u.Name
	return c.JSON(http.StatusOK, pf)
}

func DeleteProjectByID(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	idint, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "ID必须是整数",
			})
	}

	p, err := models.GetProjectByID(uint(idint))
	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "project不存在",
		})
	}

	t, _ := models.GetTeamByID(p.TeamID)
	if t == nil {
		log.WithFields(log.Fields{
			"project": *p,
		}).Error("project's team is not exist")

		return c.NoContent(http.StatusInternalServerError)
	}

	teamname := t.Name
	username := tokenInfo.Name

	flag := models.IsTeamMaintainer(teamname, username)
	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusForbidden,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	if err = models.DeleteProjectByID(uint(idint)); err != nil {
		log.WithFields(log.Fields{
			"project":  *p,
			"operator": tokenInfo.Name,
		}).Error("delete project fail")
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"project":  *p,
		"operator": tokenInfo.Name,
	}).Info("delete project success")

	return c.NoContent(http.StatusNoContent)
}

func MigrateProjectByID(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	idint, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "ID必须是整数",
			})
	}

	pmf := new(struct {
		DestTeam string `json:"dest_team" validate:"required,max=20"`
	})

	//pf := new(ProjectForm)

	if err := c.Bind(pmf); err != nil {
		//fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if err := c.Validate(pmf); err != nil {
		//fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	p, err := models.GetProjectByID(uint(idint))
	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "project不存在",
		})
	}

	t, _ := models.GetTeamByID(p.TeamID)
	if t == nil {
		log.WithFields(log.Fields{
			"project": *p,
		}).Error("project's team is not exist")

		return c.NoContent(http.StatusInternalServerError)
	}

	teamname := t.Name
	username := tokenInfo.Name

	flag := models.IsTeamMaintainer(teamname, username)
	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusForbidden,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	to, _ := models.GetTeamByName(pmf.DestTeam)
	if to == nil {
		return c.NoContent(http.StatusNotFound)
	}

	teamname = to.Name

	flag = models.IsTeamMaintainer(teamname, username)
	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusForbidden,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	if err := models.MigrateProjectByID(p.ID, to.ID); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"project":  *p,
		"operator": tokenInfo.Name,
	}).Info("migrate project success")

	return c.NoContent(http.StatusOK)
}

func GetProjectApis(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	p_id := c.Param("id")
	p_idstr, _ := strconv.Atoi(p_id)

	username := tokenInfo.Name

	p, _ := models.GetProjectByID(uint(p_idstr))
	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "项目不存在",
		})
	}

	t, _ := models.GetTeamByID(p.TeamID)
	if t == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "team不存在",
		})
	}

	//username := tokenInfo.Name
	teamname := t.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag {
		flag = models.IsTeamMember(teamname, username)
	}

	if !flag {
		flag = models.IsTeamReader(teamname, username)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusForbidden,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	apis, _ := models.GetProjectApis(p.ID)
	return c.JSON(http.StatusOK, apis)
}

func GetProjectApiGroups(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	p_id := c.Param("id")
	p_idstr, _ := strconv.Atoi(p_id)

	username := tokenInfo.Name

	p, _ := models.GetProjectByID(uint(p_idstr))
	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "项目不存在",
		})
	}

	t, _ := models.GetTeamByID(p.TeamID)
	if t == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "team不存在",
		})
	}

	//username := tokenInfo.Name
	teamname := t.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag {
		flag = models.IsTeamMember(teamname, username)
	}

	if !flag {
		flag = models.IsTeamReader(teamname, username)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusForbidden,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	api_groups, _ := models.GetProjectApiGroups(p.ID)
	return c.JSON(http.StatusOK, api_groups)
}
