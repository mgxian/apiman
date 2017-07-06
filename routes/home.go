package routes

import (
	"net/http"

	"github.com/labstack/echo"
	//"github.com/labstack/echo/middleware"
	"github.com/will835559313/apiman/models"
)

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "index")
}

func Home(c echo.Context) error {
	//u := models.User{ID: 123456, Name: "will", Nickname: "毛广献", Password: "md5"}
	u := models.User{Name: "will", Nickname: "毛广献", Password: "md5"}
	nickname := u.GetMyName()
	return c.JSON(http.StatusOK, echo.Map{"user": nickname})
}
