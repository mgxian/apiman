package apigroup

import (
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/models"
	"github.com/will835559313/apiman/pkg/jwt"
	//"gopkg.in/go-playground/validator.v9"
)

type ApiGroupForm struct {
	ID          uint   `json:"id"`
	Name        string `json:"name" validate:"required,max=20"`
	Description string `json:"description" validate:"max=100"`
	Creator     string `json:"creator"`
	Project     uint   `json:"project_id"`
}

func CreateApiGroup(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	intstr, _ := strconv.Atoi(id)
	username := tokenInfo.Name

	p, _ := models.GetProjectByID(uint(intstr))
	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "project不存在",
		})
	}

	t, _ := models.GetTeamByID(p.TeamID)
	if t == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "team不存在",
		})
	}

	u, _ := models.GetUserByName(username)
	if u == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "用户不存在",
		})
	}

	teamname := t.Name
	//projectname := p.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag {
		flag = models.IsTeamMember(teamname, username)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	ag := new(models.ApiGroup)
	agf := new(ApiGroupForm)

	if err := c.Bind(agf); err != nil {
		//fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if err := c.Validate(agf); err != nil {
		//fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	agf.ID = 0
	copier.Copy(ag, agf)
	ag.CreatorID = u.ID
	ag.ProjectID = p.ID

	if err := models.CreateApiGroup(ag); err != nil {
		log.WithFields(log.Fields{
			"apigroup": *ag,
			"operator": username,
		}).Error("create apigroup fail")
		return c.NoContent(http.StatusInternalServerError)
	}

	agf.ID = ag.ID
	agf.Creator = u.Name
	agf.Project = ag.ProjectID

	log.WithFields(log.Fields{
		"apigroup": *ag,
		"operator": username,
	}).Info("create apigroup success")

	return c.JSON(http.StatusCreated, agf)
}

func GetApiGroupByID(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	intstr, _ := strconv.Atoi(id)

	a, _ := models.GetApiGroupByID(uint(intstr))
	if a == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "api group不存在",
		})
	}

	p, _ := models.GetProjectByID(a.ProjectID)
	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "project不存在",
		})
	}

	t, _ := models.GetTeamByID(p.TeamID)
	if t == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "team不存在",
		})
	}

	u, _ := models.GetUserByID(a.CreatorID)
	if u == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "user不存在",
		})
	}

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
				"message": "你没有此权限",
			})
	}

	agf := new(ApiGroupForm)
	copier.Copy(agf, a)
	agf.Creator = u.Name
	agf.Project = p.ID

	return c.JSON(http.StatusOK, agf)
}

func UpdateApiGroupByID(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	intstr, _ := strconv.Atoi(id)
	username := tokenInfo.Name

	ag, _ := models.GetApiGroupByID(uint(intstr))
	if ag == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "api group不存在",
		})
	}

	p, _ := models.GetProjectByID(ag.ProjectID)
	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "project不存在",
		})
	}

	t, _ := models.GetTeamByID(p.TeamID)
	if t == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "team不存在",
		})
	}

	u, _ := models.GetUserByID(ag.CreatorID)
	if u == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "user不存在",
		})
	}

	teamname := t.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag {
		flag = models.IsTeamMember(teamname, username)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	aguf := new(struct {
		Name        string `json:"name" validate:"required,max=20"`
		Description string `json:"description" validate:"max=100"`
	})

	if err := c.Bind(aguf); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if err := c.Validate(aguf); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	copier.Copy(ag, aguf)

	if err := models.UpdateApiGroup(ag); err != nil {
		log.WithFields(log.Fields{
			"apigroup": *ag,
			"operator": username,
		}).Error("update apigroup info fail")
		return c.NoContent(http.StatusInternalServerError)
	}

	agf := new(ApiGroupForm)
	copier.Copy(agf, ag)
	agf.Creator = u.Name
	agf.Project = p.ID

	log.WithFields(log.Fields{
		"apigroup": *ag,
		"operator": username,
	}).Error("update apigroup success")

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

	//fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	intstr, _ := strconv.Atoi(id)

	ag, _ := models.GetApiGroupByID(uint(intstr))
	if ag == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "api group不存在",
		})
	}

	p, _ := models.GetProjectByID(ag.ProjectID)
	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "project不存在",
		})
	}

	t, _ := models.GetTeamByID(p.TeamID)
	if t == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "team不存在",
		})
	}

	username := tokenInfo.Name
	teamname := t.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag {
		flag = models.IsTeamMember(teamname, username)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	if err := models.DeleteApiGroupByID(uint(intstr)); err != nil {
		log.WithFields(log.Fields{
			"apigroup": *ag,
			"operator": username,
		}).Error("delete apigroup fail")
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"apigroup": *ag,
		"operator": username,
	}).Error("delete apigroup success")
	return c.NoContent(http.StatusNoContent)
}

func GetApiGroupApis(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	intstr, _ := strconv.Atoi(id)

	a, _ := models.GetApiGroupByID(uint(intstr))
	if a == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "api group不存在",
		})
	}

	p, _ := models.GetProjectByID(a.ProjectID)
	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "project不存在",
		})
	}

	t, _ := models.GetTeamByID(p.TeamID)
	if t == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "team不存在",
		})
	}

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
				"message": "你没有此权限",
			})
	}

	apis, _ := models.GetApiGroupApis(a.ID)

	return c.JSON(http.StatusOK, apis)
}
