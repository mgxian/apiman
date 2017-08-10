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
	"github.com/will835559313/apiman/pkg/mail"
)

type UserForm struct {
	Name      string `json:"name" validate:"required,max=20"`
	Nickname  string `json:"nickname" validate:"required,max=20"`
	Email     string `json:"email" validate:"required,max=50,email"`
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
		return c.JSON(http.StatusForbidden,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	oldUser, err := models.GetUserByName(name)
	if oldUser == nil {
		return c.NoContent(http.StatusNotFound)
	}

	newUser := new(struct {
		Nickname  string `json:"nickname" validate:"min=1,max=20"`
		Email     string `json:"email" validate:"min=3,max=50,email"`
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

	if newUser.Nickname != "" {
		oldUser.Nickname = newUser.Nickname
	}
	if newUser.AvatarUrl != "" {
		oldUser.AvatarUrl = newUser.AvatarUrl
	}
	if newUser.Email != "" {
		oldUser.Email = newUser.Email
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
		return c.JSON(http.StatusForbidden,
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
		return c.JSON(http.StatusForbidden,
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

type SendUserResetPasswordLinkForm struct {
	Email string `json:"email" validate:"min=3,max=50,email"`
}

func SendUserResetPasswordLink(c echo.Context) error {
	username := c.Param("username")
	sendUserResetPasswordLink := new(SendUserResetPasswordLinkForm)

	if err := c.Bind(sendUserResetPasswordLink); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if err := c.Validate(sendUserResetPasswordLink); err != nil {
		//fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	u, _ := models.GetUserByName(username)
	if u == nil {
		return c.NoContent(http.StatusNotFound)
	}

	if u.Email != sendUserResetPasswordLink.Email {
		return c.JSON(http.StatusConflict, echo.Map{
			"message": "用户名与邮箱不匹配",
		})
	}

	var token string
	var err error
	if token, err = jwt.GetToken(u.Name, false); err != nil {
		c.NoContent(http.StatusInternalServerError)
	}

	resetLink := "http://" + c.Request().Host + "/users/" + username +
		"/reset_password?token=" + token

	//fmt.Println("---", resetLink)

	content := "你好 <b>" + u.Nickname +
		`</b>!<br>有人请求修改你的密码，你可以点击下面的链接去操作。<br>` +
		`<a href="` + resetLink + `">修改密码</a>` +
		`<br>如果你自己没有请求修改，请忽略本邮件。<br>
		一旦你访问了该链接将会创建新的密码；不访问链接密码不会被修改。`

	//fmt.Println(content)

	if err := mail.SendHtml([]string{u.Email}, "apiman - 重置密码", content); err != nil {
		fmt.Println(err)
	}

	return c.NoContent(http.StatusOK)
}

func ResetUserPassword(c echo.Context) error {
	username := c.Param("username")
	resetToken := c.QueryParam("token")
	if resetToken == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	tokenInfo, err := jwt.ParseToken(resetToken)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "token不正确",
		})
	}

	if tokenInfo.Name != username {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "用户名与token不匹配",
		})
	}

	password := new(struct {
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

	u, err := models.GetUserByName(username)
	if u == nil {
		return c.NoContent(http.StatusNotFound)
	}

	if err := models.SetPassword(u.ID, password.NewPassword); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

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

func GetUserProjects(c echo.Context) error {
	_, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	username := c.Param("username")

	u, _ := models.GetUserByName(username)
	if u == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "用户不存在",
		})
	}

	ps, _ := models.GetUserProjects(u.ID)

	return c.JSON(http.StatusOK, ps)
}
