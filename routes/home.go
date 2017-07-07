package routes

import (
	"net/http"

	"github.com/labstack/echo"
	//"github.com/labstack/echo/middleware"
	//"github.com/will835559313/apiman/models"
)

func Index(c echo.Context) error {
	return c.String(http.StatusOK, "index")
}
