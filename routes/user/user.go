package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/models"
	"github.com/will835559313/apiman/pkg/jwt"
)

type UserForm struct {
	Name      string `json:"name" validate:"required"`
	Nickname  string `json:"nickname" validate:"required,max=20"`
	Password  string `json:"password" validate:"required,min=6"`
	AvatarUrl string `json:"avatar_url"`
}

func CreateUser(c echo.Context) (err error) {
	uf := new(UserForm)
	if err := c.Bind(uf); err != nil {
		//fmt.Println(err)
		log.Error("user bind error")
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if err := c.Validate(uf); err != nil {
		//fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	u, err := models.GetUserByName(uf.Name)

	if u != nil {
		return c.NoContent(http.StatusConflict)
	}

	u = new(models.User)
	copier.Copy(u, uf)

	if err := models.CreateUser(u); err != nil {
		//log.Info("user create error")
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"user": *u,
	}).Info("user register success")

	//return c.JSONPretty(http.StatusCreated, u, "  ")
	return c.JSON(http.StatusCreated, u)
}

func GetUserByID(c echo.Context) (err error) {
	_, err = jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	idstr := c.Param("id")
	idint, err := strconv.Atoi(idstr)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "ID必须是整数",
			})
	}

	u, _ := models.GetUserByID(uint(idint))
	if u == nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, u)
}

func GetUserByName(c echo.Context) (err error) {
	_, err = jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	name := c.Param("username")
	u := new(models.User)
	u, err = models.GetUserByName(name)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, u)
}

func UpdateUserByName(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	name := c.Param("username")

	if name != tokenInfo.Name && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	oldUser, err := models.GetUserByName(name)
	if oldUser == nil {
		return c.NoContent(http.StatusNotFound)
	}

	newUser := new(struct {
		Nickname  string `json:"nickname" validate:"required,max=20"`
		AvatarUrl string `json:"avatar_url"`
	})

	err = c.Bind(newUser)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if err = c.Validate(newUser); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	oldUser.Nickname = newUser.Nickname
	if newUser.AvatarUrl != "" {
		oldUser.AvatarUrl = newUser.AvatarUrl
	}

	err = models.UpdateUser(oldUser)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	u, err := models.GetUserByName(name)
	if u == nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"user":     *u,
		"operator": tokenInfo.Name,
	}).Info("user info update success")

	return c.JSON(http.StatusOK, u)
}

func DeleteUserByName(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	name := c.Param("username")

	if name != tokenInfo.Name && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	err = models.DeleteUserByUsername(name)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	log.WithFields(log.Fields{
		"username": name,
		"operator": tokenInfo.Name,
	}).Info("delete user success")

	return c.NoContent(http.StatusNoContent)
}

func ChangeUserPassword(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	name := c.Param("username")

	if name != tokenInfo.Name && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	password := new(struct {
		OldPassword string `json:"old_password" validate:"required,min=6"`
		NewPassword string `json:"new_password" validate:"required,min=6`
	})

	if err := c.Bind(password); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if err = c.Validate(password); err != nil {
		//fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	u, err := models.GetUserByName(name)
	if u == nil {
		return c.NoContent(http.StatusNotFound)
	}

	err = models.ChangeUserPassword(password.OldPassword, password.NewPassword, u.ID)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "密码不正确",
		})
	}

	log.WithFields(log.Fields{
		"username": name,
		"operator": tokenInfo.Name,
	}).Info("change user password success")

	return c.NoContent(http.StatusOK)
}

func GetToken(c echo.Context) error {
	login := new(struct {
		Name     string `json:"name" validate:"required,max=20"`
		Password string `json:"password" validate:"required,min=6"`
	})

	if err := c.Bind(login); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if err := c.Validate(login); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	token, err := models.GetToken(login.Name, login.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "用户名或密码错误",
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"access_token": token,
	})
}
