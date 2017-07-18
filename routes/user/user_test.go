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
	userJSON            = `{"name":"will","nickname":"毛广献","password":"mgx123","avatar_url":"http://ojz1mcltu.bkt.clouddn.com/animals-august2015.jpg"}`
	badUserJSON         = `{"name":"will","nickname":"毛广献","password":"mgx"}`
	newUserJSON         = `{"name":"will","nickname":"毛","password":"mgx123","avatar_url":"http://ojz1mcltu.bkt.clouddn.com/animals-august2015.jpg"}`
	newBadUserJSON      = `{"name":"will","nickname":"111111111122222222223"}`
	authJSON            = `{"name":"will", "password":"mgx123"}`
	badAuthJSON         = `{"name":"will", "password":"mgx"}`
	restPasswordJson    = `{"old_password":"mgx123","new_password":"will123"}`
	badRestPasswordJson = `{"old_password":"qwerty","new_password":"will123"}`
	access_token        = ""
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
	// models.DB.DropTableIfExists(&models.User{})
	models.DbMigrate()

	// set jwt
	jwt.JwtInint()
}

// normal request
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

// bad request
func TestBadCreateUser(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/users", strings.NewReader(badUserJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, CreateUser(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

// normal request
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

// bad request
func TestBadGetToken(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/oauth2/token", strings.NewReader(badAuthJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, GetToken(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

// normal request
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

// bad request
func TestBadGetUser(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	//req.Header.Set(echo.HeaderAuthorization, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:username")
	c.SetParamNames("username")
	c.SetParamValues("willmgx")

	// Assertions
	if assert.NoError(t, GetUserByName(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

// normal request
func TestUpdateUser(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
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

// bad request
func TestBadUpdateUser(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(newBadUserJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/users/:username")
	c.SetParamNames("username")
	c.SetParamValues("will")

	// Assertions
	if assert.NoError(t, UpdateUserByName(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

// normal request
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

// bad request
func TestBadChangeUserPassword(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/users/:username/change_password", strings.NewReader(badRestPasswordJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("username")
	c.SetParamValues("will")

	// Assertions
	if assert.NoError(t, ChangeUserPassword(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}
}

// normal request
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

	// delete user table
	// models.DB.DropTableIfExists(&models.User{})
}

// bad request
func TestBadDeleteUser(t *testing.T) {
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
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}

	// delete user table
	// models.DB.DropTableIfExists(&models.User{})
}
