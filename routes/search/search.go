package search

import (
	//"fmt"
	"net/http"
	"strings"

	"github.com/jinzhu/copier"
	//"github.com/bitly/go-simplejson"
	"github.com/labstack/echo"
	//log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/models"
	"github.com/will835559313/apiman/pkg/jwt"
	"github.com/will835559313/apiman/routes/api"
	"github.com/will835559313/apiman/routes/apigroup"
	"github.com/will835559313/apiman/routes/project"
	"github.com/will835559313/apiman/routes/team"
)

func Search(c echo.Context) error {
	_, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	search_type := c.Param("type")
	q := c.QueryParam("q")
	sort := c.QueryParam("sort")
	order := c.QueryParam("order")

	//fmt.Println(search_type, q, strings.Title(sort), strings.ToLower(order))
	ids, _ := models.Search(q, search_type, strings.Title(sort), strings.ToLower(order))

	switch search_type {
	case "user":
		users := make([]*models.User, 0)
		for _, id := range ids {
			u, _ := models.GetUserByID(id)
			if u != nil {
				users = append(users, u)
			}
		}

		return c.JSON(http.StatusOK, echo.Map{
			"users": users,
		})
	case "team":
		teams := make([]*team.TeamForm, 0)
		for _, id := range ids {
			t, _ := models.GetTeamByID(id)
			if t != nil {
				tf := new(team.TeamForm)
				copier.Copy(tf, t)
				u, _ := models.GetUserByID(t.CreatorID)
				tf.Creator = u.Name
				teams = append(teams, tf)
			}
		}

		return c.JSON(http.StatusOK, echo.Map{
			"teams": teams,
		})
	case "project":
		projects := make([]*project.ProjectForm, 0)
		for _, id := range ids {
			p, _ := models.GetProjectByID(id)
			if p != nil {
				pf := new(project.ProjectForm)
				copier.Copy(pf, p)
				u, _ := models.GetUserByID(p.CreatorID)
				pf.Creator = u.Name
				t, _ := models.GetTeamByID(p.TeamID)
				pf.Team = t.Name
				projects = append(projects, pf)
			}
		}

		return c.JSON(http.StatusOK, echo.Map{
			"projects": projects,
		})
	case "api_group":
		api_groups := make([]*apigroup.ApiGroupForm, 0)
		for _, id := range ids {
			ag, _ := models.GetApiGroupByID(id)
			if ag != nil {
				agf := new(apigroup.ApiGroupForm)
				copier.Copy(agf, ag)
				u, _ := models.GetUserByID(ag.CreatorID)
				//p, _ := models.GetProjectByID(ag.ProjectID)
				agf.Creator = u.Name
				agf.Project = ag.ProjectID
				api_groups = append(api_groups, agf)
			}
		}

		return c.JSON(http.StatusOK, echo.Map{
			"api_groups": api_groups,
		})
	case "api":
		apis := make([]*api.ApiBaseInfo, 0)
		for _, id := range ids {
			a, _ := models.GetApiByID(id)
			if a != nil {
				apif := new(api.ApiBaseInfo)
				copier.Copy(apif, a)
				u, _ := models.GetUserByID(a.CreatorID)
				apif.Creator = u.Name
				apis = append(apis, apif)
			}
		}

		return c.JSON(http.StatusOK, echo.Map{
			"apis": apis,
		})
	default:
		return c.NoContent(http.StatusBadRequest)
	}
}
