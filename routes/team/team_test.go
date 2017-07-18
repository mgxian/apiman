package team

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
	"github.com/will835559313/apiman/routes/user"
	"gopkg.in/go-playground/validator.v9"
)

var (
	teamJSON     = `{"name":"famulei","description":"team","creator":"will","avatar_url":"http://www.famulei.com/images/index_v4/slogan.png"}`
	newTeamJSON  = `{"name":"famulei","description":"a great team","avatar_url":"http://www.famulei.com/images/index_v3/slogan.png"}`
	authJSON     = `{"name":"will", "password":"mgx123"}`
	access_token = ""

	badTeamJSON    = `{"name":"famulei111111111111122222222222222222222","description":"team","creator":"will"}`
	badNewTeamJSON = `{"name":"famulei","description":"a great team"}`

	userJSON = `{"name":"will","nickname":"毛广献","password":"mgx123","avatar_url":"http://ojz1mcltu.bkt.clouddn.com/animals-august2015.jpg"}`
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

func createUser() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	user.CreateUser(c)
}

func TestGetToken(t *testing.T) {
	// create user
	createUser()

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
	fmt.Println(access_token)
}

func TestCreateTeam(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/teams", strings.NewReader(teamJSON))
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, CreateTeam(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		t1 := new(models.Team)
		t2 := new(models.Team)
		json.Unmarshal([]byte(rec.Body.String()), t2)
		json.Unmarshal([]byte(teamJSON), t1)
		assert.Equal(t, t1.Name, t2.Name)
		assert.Equal(t, t1.Description, t2.Description)
	}
}

// bad request
func TestBadCreateTeam(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/teams", strings.NewReader(badTeamJSON))
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, CreateTeam(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	}
}

func TestGetTeam(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/teams/:teamname", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("teamname")
	c.SetParamValues("famulei")

	// Assertions
	if assert.NoError(t, GetTeamByName(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		t1 := new(models.Team)
		t2 := new(models.Team)
		json.Unmarshal([]byte(rec.Body.String()), t2)
		json.Unmarshal([]byte(teamJSON), t1)
		assert.Equal(t, t1.Name, t2.Name)
		assert.Equal(t, t1.Description, t2.Description)
	}
}

// bad request
func TestBadGetTeam(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/teams/:teamname", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("teamname")
	c.SetParamValues("weoriqwp")

	// Assertions
	if assert.NoError(t, GetTeamByName(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
	}
}

func TestUpdateTeam(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(newTeamJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/teams/:teamname")
	c.SetParamNames("teamname")
	c.SetParamValues("famulei")

	// Assertions
	if assert.NoError(t, UpdateTeamByName(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		t1 := new(models.Team)
		t2 := new(models.Team)
		json.Unmarshal([]byte(rec.Body.String()), t2)
		json.Unmarshal([]byte(newTeamJSON), t1)
		assert.Equal(t, t1.Name, t2.Name)
		assert.Equal(t, t1.Description, t2.Description)
	}
}

// bad request
func TestBadUpdateTeam(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(badNewTeamJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	//req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/teams/:teamname")
	c.SetParamNames("teamname")
	c.SetParamValues("famulei")

	// Assertions
	if assert.NoError(t, UpdateTeamByName(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}
}

func TestDeleteTeam(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/teams/:teamname", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("teamname")
	c.SetParamValues("famulei")

	// Assertions
	if assert.NoError(t, DeleteTeamByName(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}

	// delete team table
	// models.DB.DropTableIfExists(&models.Team{})
}

// bad request
func TestBadDeleteTeam(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/teams/:teamname", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("teamname")
	c.SetParamValues("famulei")

	// Assertions
	if assert.NoError(t, DeleteTeamByName(c)) {
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	}

	// delete team table
	// models.DB.DropTableIfExists(&models.Team{})
}
