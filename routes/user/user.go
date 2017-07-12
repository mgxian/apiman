package user

import (
	"fmt"
	//"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/models"
	"github.com/will835559313/apiman/pkg/jwt"
	//"gopkg.in/go-playground/validator.v9"
)

type UserForm struct {
	Name      string `json:"name" validate:"required"`
	Nickname  string `json:"nickname" validate:"required,max=20"`
	Password  string `json:"password" validate:"required,min=6"`
	AvatarUrl string `json:"avatar_url"`
}

func CreateUser(c echo.Context) (err error) {
	uf := new(UserForm)
	if err = c.Bind(uf); err != nil {
		fmt.Println(err)
		log.Info("user bind error")
		return
	}

	if err = c.Validate(uf); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "request data is not right",
		})
	}

	u, err := models.GetUserByName(uf.Name)

	if u != nil {
		return c.NoContent(http.StatusConflict)
	}

	fmt.Println("build finish")
	if err = c.Validate(uf); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("validator finish")

	u = new(models.User)
	u.Name = uf.Name
	u.Nickname = uf.Nickname
	u.Password = uf.Password
	u.AvatarUrl = uf.AvatarUrl
	if err := models.CreateUser(u); err != nil {
		log.Info("user create error")
		return c.NoContent(http.StatusInternalServerError)
	}
	log.WithFields(log.Fields{
		"username": u.Name,
	}).Info("user create username=" + u.Name)
	//return c.JSONPretty(http.StatusCreated, u, "  ")
	return c.JSON(http.StatusCreated, u)
}

func GetUserByID(c echo.Context) (err error) {
	idstr := c.Param("id")
	idint, err := strconv.Atoi(idstr)
	u, err := models.GetUserByID(uint(idint))
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}
	if err = c.Validate(u); err != nil {
		log.Info("validator error")
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSONPretty(http.StatusOK, u, "  ")
}

func GetUserByName(c echo.Context) (err error) {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	name := c.Param("username")
	u := new(models.User)
	u, err = models.GetUserByName(name)
	if err != nil {
		return c.NoContent(http.StatusNotFound)
	}

	//	if err = c.Validate(u); err != nil {
	//		log.Info("validator error")
	//		fmt.Println(err)
	//		//return c.NoContent(http.StatusNotFound)
	//		return err
	//	}
	return c.JSONPretty(http.StatusOK, u, "  ")
}

func UpdateUserByName(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	fmt.Println(tokenInfo.Name)

	name := c.Param("username")

	if name != tokenInfo.Name && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	oldUser, err := models.GetUserByName(name)
	if oldUser == nil {
		return c.NoContent(http.StatusNotFound)
	}

	//	type UserUpdateForm struct {
	//	Nickname  string `json:"nickname" validate:"required,max=20"`
	//	AvatarUrl string `json:"avatar_url"`
	//	}

	newUser := new(struct {
		Nickname  string `json:"nickname" validate:"required,max=20"`
		AvatarUrl string `json:"avatar_url"`
	})
	c.Bind(newUser)

	if err = c.Validate(newUser); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "request data is not right",
		})
	}

	oldUser.Nickname = newUser.Nickname
	oldUser.AvatarUrl = newUser.AvatarUrl

	err = models.UpdateUser(oldUser)

	u, err := models.GetUserByName(name)
	if u == nil {
		log.WithFields(log.Fields{
			"username": name,
		}).Info(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

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

	fmt.Println(tokenInfo.Name)

	name := c.Param("username")

	if name != tokenInfo.Name && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	err = models.DeleteUserByUsername(name)
	if err != nil {
		log.Error("no such user")
	}
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

	fmt.Println(tokenInfo.Name)

	name := c.Param("username")

	if name != tokenInfo.Name && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "you have not this permisson",
			})
	}

	password := new(struct {
		OldPassword string `json:"old_password" validate:"required,min=6"`
		NewPassword string `json:"new_password" validate:"required,min=6`
	})

	if err := c.Bind(password); err != nil {
		fmt.Println(err)
	}

	if err = c.Validate(password); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "request data is not right",
		})
	}

	u, err := models.GetUserByName(name)
	if u == nil {
		return c.NoContent(http.StatusNotFound)
	}
	err = models.ChangeUserPassword(password.OldPassword, password.NewPassword, u.ID)
	if err != nil {
		//data := map[string]string{}
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": string(err.Error()),
		})
	}
	return c.NoContent(http.StatusOK)
}

func GetToken(c echo.Context) error {
	login := new(struct {
		Name     string `json:"name" validate:"required,max=20"`
		Password string `json:"password" validate:"required,min=6"`
	})

	if err := c.Bind(login); err != nil {
		log.WithFields(log.Fields{
			"username": login.Name,
		}).Info("login fail")
		return c.JSONPretty(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		}, "  ")
	}

	if err := c.Validate(login); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "request data is not right",
		})
	}

	if token, err := models.GetToken(login.Name, login.Password); err != nil {
		return c.JSONPretty(http.StatusUnauthorized, echo.Map{
			"message": "用户名或密码错误",
		}, "  ")
	} else {
		return c.JSONPretty(http.StatusCreated, echo.Map{
			"access_token": token,
		}, "  ")
	}
}
