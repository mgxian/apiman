package apigroup

import (
	"encoding/json"
	"errors"
	//"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/bitly/go-simplejson"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/will835559313/apiman/models"
	"github.com/will835559313/apiman/pkg/jwt"
	"github.com/will835559313/apiman/pkg/log"
	"github.com/will835559313/apiman/pkg/setting"
	"github.com/will835559313/apiman/routes/user"
	"gopkg.in/go-playground/validator.v9"
)

var (
	apigroupJSON    = `{"name":"becai","description":"菠菜"}`
	newApiGroupJSON = `{"name":"bocai","description":"博彩"}`
	authJSON        = `{"name":"will", "password":"mgx123"}`
	access_token    = ""
	apigroupID      = uint(0)
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
	// models.DB.DropTableIfExists(&models.Team{})
	models.DbMigrate()

	// set jwt
	jwt.JwtInint()
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
	if assert.NoError(t, user.GetToken(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), "access_token")
	}

	// save token
	js, err := simplejson.NewJson([]byte(rec.Body.String()))
	if err != nil {
		assert.Error(t, errors.New("save token error"))
	}
	access_token, _ = js.Get("access_token").String()
	//fmt.Println(access_token)
}

func TestCreateApiGroup(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/projects/:id/apigroups", strings.NewReader(apigroupJSON))
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("2")

	// Assertions
	if assert.NoError(t, CreateApiGroup(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		t1 := new(models.ApiGroup)
		t2 := new(models.ApiGroup)
		json.Unmarshal([]byte(rec.Body.String()), t2)
		json.Unmarshal([]byte(apigroupJSON), t1)
		assert.Equal(t, t1.Name, t2.Name)
		assert.Equal(t, t1.Description, t2.Description)
		apigroupID = t2.ID
	}
}

func TestGetApiGroup(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/apigroups/:id", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	idstr := strconv.Itoa(int(apigroupID))
	c.SetParamValues(idstr)

	// Assertions
	if assert.NoError(t, GetApiGroupByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		t1 := new(models.ApiGroup)
		t2 := new(models.ApiGroup)
		json.Unmarshal([]byte(rec.Body.String()), t2)
		json.Unmarshal([]byte(apigroupJSON), t1)
		assert.Equal(t, t1.Name, t2.Name)
		assert.Equal(t, t1.Description, t2.Description)
	}
}

func TestUpdateApiGroup(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(newApiGroupJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/apigroups/:id")
	idstr := strconv.Itoa(int(apigroupID))
	c.SetParamNames("id")
	c.SetParamValues(idstr)

	// Assertions
	if assert.NoError(t, UpdateApiGroupByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		t1 := new(models.ApiGroup)
		t2 := new(models.ApiGroup)
		json.Unmarshal([]byte(rec.Body.String()), t2)
		json.Unmarshal([]byte(newApiGroupJSON), t1)
		assert.Equal(t, t1.Name, t2.Name)
		assert.Equal(t, t1.Description, t2.Description)
	}
}

func TestDeleteApiGroup(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/apigroups/:id", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	idstr := strconv.Itoa(int(apigroupID))
	c.SetParamNames("id")
	c.SetParamValues(idstr)

	// Assertions
	if assert.NoError(t, DeleteApiGroupByID(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}

	// delete team table
	// models.DB.DropTableIfExists(&models.Team{})
}
