package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bitly/go-simplejson"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/will835559313/apiman/models"
	"github.com/will835559313/apiman/pkg/jwt"
	"github.com/will835559313/apiman/pkg/log"
	"github.com/will835559313/apiman/pkg/setting"
	"gopkg.in/go-playground/validator.v9"
)

var (
	userJSON         = `{"name":"will","nickname":"毛广献","password":"mgx123","avatar_url":"http://ojz1mcltu.bkt.clouddn.com/animals-august2015.jpg"}`
	newUserJSON      = `{"name":"will","nickname":"毛","password":"mgx123","avatar_url":"http://ojz1mcltu.bkt.clouddn.com/animals-august2015.jpg"}`
	authJSON         = `{"name":"will", "password":"mgx123"}`
	restPasswordJson = `{"old_password":"mgx123","new_password":"will123"}`
	access_token     = ""
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func init() {
	// load config
	setting.Conf = "C:/Users/will8/go/src/github.com/will835559313/apiman/conf/app.conf"
	setting.CustomConf = "C:/Users/will8/go/src/github.com/will835559313/apiman/conf/custom.app.conf"
	setting.NewConfig()

	// set logger
	log.LoggerInit()

	// get db connection
	models.Dbinit()

	// migrate tables
	models.DB.DropTableIfExists(&models.User{})
	models.DbMigrate()

	// set jwt
	jwt.JwtInint()
}

func TestCreateUser(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, CreateUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		u1 := new(models.User)
		u2 := new(models.User)
		json.Unmarshal([]byte(rec.Body.String()), u2)
		json.Unmarshal([]byte(userJSON), u1)
		assert.Equal(t, u1.Name, u2.Name)
		assert.Equal(t, u1.Nickname, u2.Nickname)
	}
}

func TestGetToken(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/oauth2/token", strings.NewReader(authJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, GetToken(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "access_token")
	}

	// save token
	js, err := simplejson.NewJson([]byte(rec.Body.String()))
	if err != nil {
		assert.Error(t, errors.New("save token error"))
	}
	access_token, _ = js.Get("access_token").String()
	fmt.Println(access_token)
}

func TestGetUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	//req.Header.Set(echo.HeaderAuthorization, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:username")
	c.SetParamNames("username")
	c.SetParamValues("will")

	// Assertions
	if assert.NoError(t, GetUserByName(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		u1 := new(models.User)
		u2 := new(models.User)
		json.Unmarshal([]byte(rec.Body.String()), u2)
		json.Unmarshal([]byte(userJSON), u1)
		assert.Equal(t, u1.Name, u2.Name)
		assert.Equal(t, u1.Nickname, u2.Nickname)
	}
}

func TestUpdateUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(newUserJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:username")
	c.SetParamNames("username")
	c.SetParamValues("will")

	// Assertions
	if assert.NoError(t, UpdateUserByName(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		u1 := new(models.User)
		u2 := new(models.User)
		json.Unmarshal([]byte(rec.Body.String()), u2)
		json.Unmarshal([]byte(newUserJSON), u1)
		assert.Equal(t, u1.Name, u2.Name)
		assert.Equal(t, u1.Nickname, u2.Nickname)
	}
}

func TestChangeUserPassword(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/users/:username/change_password", strings.NewReader(restPasswordJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("username")
	c.SetParamValues("will")

	// Assertions
	if assert.NoError(t, ChangeUserPassword(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestDeleteUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/users/:username", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	//req.Header.Set(echo.HeaderAuthorization, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("username")
	c.SetParamValues("will")

	// Assertions
	if assert.NoError(t, DeleteUserByName(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}
