package api

import (
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/will835559313/apiman/routes/apigroup"
	"github.com/will835559313/apiman/routes/project"
	"github.com/will835559313/apiman/routes/team"
	"github.com/will835559313/apiman/routes/user"
	"gopkg.in/go-playground/validator.v9"
)

var (
	apiBaseJSON = `{"base_info":{"id":1,"name":"用户","description":"","creator":"will","project":0,"group":0,"uri":"/users","protocol":1,"method":1}}`
	apiFullJSON = `{"base_info":{"id":1,"name":"用户","description":"","creator":"will","project":0,"group":0,"uri":"/users","protocol":1,"method":1},"request":{"headers":[{"id":1,"name":"Authorization","value":"Bearer eyJhbGciOiJI","description":"认证信息","api_id":1},{"id":2,"name":"Content-Type","value":"application/json","description":"数据响应的格式","remark":"","api_id":1}],"parameters":[{"id":1,"name":"users","value":"","type":7,"required":true,"description":"","remark":"用户集合","api_id":1,"parent_id":0,"sub_parameter":[{"id":2,"name":"username","value":"","type":2,"required":true,"description":"","remark":"用户名","api_id":1,"parent_id":1,"sub_parameter":[{"id":2,"name":"firt_name","value":"","type":2,"required":true,"description":"","remark":"用户姓氏","api_id":1,"parent_id":1,"sub_parameter":null},{"id":2,"name":"second_name","value":"","type":2,"description":"","remark":"用户名字","api_id":1,"parent_id":1,"sub_parameter":null}]},{"id":2,"name":"id","value":"","type":2,"description":"","remark":"用户id","api_id":1,"parent_id":1,"sub_parameter":null}]}]},"response":{"headers":[{"id":2,"name":"Content-Type","value":"application/json","description":"数据响应的格式","remark":"","api_id":1}],"parameters":[{"id":1,"name":"users","value":"","type":7,"required":true,"description":"","remark":"用户集合","api_id":1,"parent_id":0,"sub_parameter":[{"id":2,"name":"username","value":"","type":2,"required":true,"description":"","remark":"用户名","api_id":1,"parent_id":1,"sub_parameter":[{"id":2,"name":"firt_name","value":"","type":2,"required":true,"description":"","remark":"用户姓氏","api_id":1,"parent_id":1,"sub_parameter":null},{"id":2,"name":"second_name","value":"","type":2,"description":"","remark":"用户名字","api_id":1,"parent_id":1,"sub_parameter":null}]},{"id":2,"name":"id","value":"","type":2,"description":"","remark":"用户id","api_id":1,"parent_id":1,"sub_parameter":null}]}]}}`

	authJSON     = `{"name":"will", "password":"mgx123"}`
	access_token = ""
	apigroupID   = uint(0)
	apiID        = uint(0)

	userJSON     = `{"name":"will","nickname":"毛广献","password":"mgx123","avatar_url":"http://ojz1mcltu.bkt.clouddn.com/animals-august2015.jpg"}`
	teamJSON     = `{"name":"famulei","description":"team","creator":"will","avatar_url":"http://www.famulei.com/images/index_v4/slogan.png"}`
	projectJSON  = `{"name":"web","description":"web版","avatar_url":"http://www.famulei.com/"}`
	apigroupJSON = `{"name":"becai","description":"菠菜"}`
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

func createTeam() {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/teams", strings.NewReader(teamJSON))
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	team.CreateTeam(c)
}

func createProject() {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/teams/:teamname/projects", strings.NewReader(projectJSON))
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("teamname")
	c.SetParamValues("famulei")

	project.CreateProject(c)
}

func createApiGroup() {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/projects/:id/apigroups", strings.NewReader(apigroupJSON))
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	apigroup.CreateApiGroup(c)
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
	//fmt.Println(access_token)
}

func TestCreateApi(t *testing.T) {
	// create team
	createTeam()

	// create project
	createProject()

	// create apigroup
	createApiGroup()

	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.POST, "/apigroups/:id/apis", strings.NewReader(apiBaseJSON))
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	// Assertions
	if assert.NoError(t, CreateApi(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		t1 := new(ApiForm)
		t2 := new(ApiForm)
		json.Unmarshal([]byte(rec.Body.String()), t2)
		json.Unmarshal([]byte(apiBaseJSON), t1)
		assert.Equal(t, t1.Name, t2.Name)
		assert.Equal(t, t1.Description, t2.Description)
		apiID = t2.ID
		fmt.Println("---------------")
		fmt.Println(apiID)
		fmt.Println("---------------")
	}
}

func TestGetApi(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/apis/:id", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	idstr := strconv.Itoa(int(apiID))
	c.SetParamValues(idstr)

	// Assertions
	if assert.NoError(t, GetApi(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		t1 := new(models.Api)
		t2 := new(models.Api)
		json.Unmarshal([]byte(rec.Body.String()), t2)
		json.Unmarshal([]byte(apiBaseJSON), t1)
		assert.Equal(t, t1.Name, t2.Name)
		assert.Equal(t, t1.Description, t2.Description)
	}
}

func TestUpdateApi(t *testing.T) {
	// Setup
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(apiFullJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/apis/:id")
	idstr := strconv.Itoa(int(apiID))
	c.SetParamNames("id")
	c.SetParamValues(idstr)

	// Assertions
	if assert.NoError(t, UpdateApi(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		t1 := new(ApiForm)
		t2 := new(ApiForm)
		json.Unmarshal([]byte(rec.Body.String()), t2)
		if err := json.Unmarshal([]byte(apiFullJSON), t1); err != nil {
			fmt.Println(err)
		}

		assert.Equal(t, t1.Name, t2.Name)
		assert.Equal(t, t1.Description, t2.Description)
		assert.NotEmpty(t, t2.Request.RequestParameters)
	}
}

func TestDeleteApi(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(echo.DELETE, "/apis/:id", nil)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+access_token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	idstr := strconv.Itoa(int(apiID))
	c.SetParamNames("id")
	c.SetParamValues(idstr)

	// Assertions
	if assert.NoError(t, DeleteApi(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}

	// delete team table
	// models.DB.DropTableIfExists(&models.Team{})
}
