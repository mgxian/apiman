package search

import (
	//"fmt"
	"net/http"
	"strings"

	//"github.com/jinzhu/copier"
	//"github.com/bitly/go-simplejson"
	"github.com/labstack/echo"
	//log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/models"
	"github.com/will835559313/apiman/pkg/jwt"
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

		return c.JSON(http.StatusOK, users)
	case "team":
		teams := make([]*models.Team, 0)
		for _, id := range ids {
			u, _ := models.GetTeamByID(id)
			if u != nil {
				teams = append(teams, u)
			}
		}

		return c.JSON(http.StatusOK, teams)
	case "project":
		projects := make([]*models.Project, 0)
		for _, id := range ids {
			u, _ := models.GetProjectByID(id)
			if u != nil {
				projects = append(projects, u)
			}
		}

		return c.JSON(http.StatusOK, projects)
	case "api_group":
		api_groups := make([]*models.ApiGroup, 0)
		for _, id := range ids {
			u, _ := models.GetApiGroupByID(id)
			if u != nil {
				api_groups = append(api_groups, u)
			}
		}

		return c.JSON(http.StatusOK, api_groups)
	case "api":
		apis := make([]*models.Api, 0)
		for _, id := range ids {
			u, _ := models.GetApiByID(id)
			if u != nil {
				apis = append(apis, u)
			}
		}

		return c.JSON(http.StatusOK, apis)
	default:
		return c.NoContent(http.StatusBadRequest)
	}
}
