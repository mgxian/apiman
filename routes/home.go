package routes

import (
	"net/http"

	"github.com/labstack/echo"
	//"github.com/labstack/echo/middleware"
)

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "index")
}

func Home(c echo.Context) error {
	return c.JSON(http.StatusOK, echo.Map{"user": "will"})
}
