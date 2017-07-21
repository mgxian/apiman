package api

import (
	"encoding/json"
	"fmt"
	//"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	//"github.com/bitly/go-simplejson"
	"github.com/labstack/echo"
	//log "github.com/sirupsen/logrus"
	"github.com/will835559313/apiman/models"
	"github.com/will835559313/apiman/pkg/jwt"
)

type RequestParameters struct {
	models.RequestParameter
	SubParameters []*RequestParameters `json:"sub_parameter"`
}

type ResponseParameters struct {
	models.ResponseParameter
	SubParameters []*ResponseParameters `json:"sub_parameter"`
}

type RequestForm struct {
	Request struct {
		RequestHeaders    []*models.RequestHeader `json:"headers"`
		RequestParameters []*RequestParameters    `json:"parameters"`
	} `json:"request"`

	Response struct {
		ResponseHeaders    []*models.ResponseHeader `json:"headers"`
		ResponseParameters []*ResponseParameters    `json:"parameters"`
	} `json:"response"`
}

type ApiForm struct {
	models.Api
	RequestForm
	Creator string `json:"creator" validate:"required,max=20"`
}

type ApiData struct {
	models.Api
	RequestForm
}

func TestApi(c echo.Context) error {
	// set base info
	api := new(ApiForm)
	api.ID = 1
	api.Name = "用户"
	api.Method = models.GET
	api.Protocol = models.HTTP
	api.URI = "/users"

	// set header
	h1 := new(models.RequestHeader)
	h1.ApiID = 1
	h1.Name = "Authorization"
	h1.Value = "Bearer eyJhbGciOiJI"
	h1.ID = 1
	h1.Description = "认证信息"

	h2 := new(models.ResponseHeader)
	h2.ApiID = 1
	h2.Name = "Content-Type"
	h2.Value = "application/json"
	h2.ID = 2
	h2.Description = "数据响应的格式"

	api.Request.RequestHeaders = make([]*models.RequestHeader, 0)
	api.Response.ResponseHeaders = make([]*models.ResponseHeader, 0)

	api.Request.RequestHeaders = append(api.Request.RequestHeaders, h1)
	api.Response.ResponseHeaders = append(api.Response.ResponseHeaders, h2)

	// set request parameter
	resp1 := new(ResponseParameters)
	resp1.ApiID = 1
	resp1.ID = 1
	resp1.Name = "users"
	resp1.ParentID = 0
	resp1.Required = true
	resp1.Type = models.ArrayObject
	resp1.Remark = "用户集合"

	resp2 := new(ResponseParameters)
	resp2.ApiID = 1
	resp2.ID = 2
	resp2.Name = "username"
	resp2.ParentID = 1
	resp2.Required = true
	resp2.Type = models.String
	resp2.Remark = "用户名"

	resp1.SubParameters = append(resp1.SubParameters, resp2)

	api.Request.RequestParameters = make([]*RequestParameters, 0)
	api.Response.ResponseParameters = append(api.Response.ResponseParameters, resp1)

	s, _ := json.Marshal(resp1)
	fmt.Printf("%s", s)
	return c.JSON(http.StatusOK, api)
}

func CreateApi(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	apif := new(ApiForm)
	if err := c.Bind(apif); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	if err := c.Validate(apif); err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "请求数据错误",
		})
	}

	groupid := c.Param("id")
	groupidstr, _ := strconv.Atoi(groupid)
	g, _ := models.GetApiGroupByID(uint(groupidstr))
	if g == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "api group不存在",
		})
	}

	username := tokenInfo.Name
	u, _ := models.GetUserByName(username)
	if u == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "用户不存在",
		})
	}

	p, _ := models.GetProjectByID(g.Project)
	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "项目不存在",
		})
	}

	userid := strconv.Itoa(int(u.ID))
	apif.Creator = userid
	apif.Group = g.ID
	apif.Project = p.ID

	if err := saveApi(apif, true); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, apif)
}

func saveApi(api *ApiForm, create bool) error {
	// save api base info
	apiBaseInfo := new(models.Api)
	copier.Copy(apiBaseInfo, api)
	if create {
		apiBaseInfo.ID = 0
	}

	userid, _ := strconv.Atoi(api.Creator)
	apiBaseInfo.Creator = uint(userid)

	if err := models.CreateApi(apiBaseInfo); err != nil {
		return err
	}

	// save header
	for _, rh := range api.Request.RequestHeaders {
		//copier.Copy(apiRequestHeader, rh)
		if create {
			rh.ID = 0
			rh.ApiID = apiBaseInfo.ID
		} else {
			rh.ApiID = 0
		}

		//fmt.Printf("%v", rh)
	}

	if err := models.BatchCreateRequestHeader(api.Request.RequestHeaders); err != nil {
		fmt.Print(err)
	}

	for _, rh := range api.Response.ResponseHeaders {
		//copier.Copy(apiResponseHeader, rh)
		if create {
			rh.ID = 0
			rh.ApiID = apiBaseInfo.ID
		} else {
			rh.ApiID = 0
		}

		//fmt.Printf("%v", rh)
	}

	if err := models.BatchCreateResponseHeader(api.Response.ResponseHeaders); err != nil {
		fmt.Print(err)
	}

	// save parameter
	if err := saveRequestParameters(api.Request.RequestParameters, uint(0), apiBaseInfo.ID); err != nil {
		fmt.Print(err)
	}

	if err := saveResponseParameters(api.Response.ResponseParameters, uint(0), apiBaseInfo.ID); err != nil {
		fmt.Print(err)
	}

	return nil
}

func saveRequestParameters(rps []*RequestParameters, p_id uint, api_id uint) error {
	if len(rps) == 0 {
		return nil
	}

	for _, rp := range rps {
		requestParameter := new(models.RequestParameter)
		copier.Copy(requestParameter, rp)
		requestParameter.ID = 0
		requestParameter.ApiID = api_id
		requestParameter.ParentID = p_id
		fmt.Printf("-------%v\n", requestParameter)
		if err := models.CreateRequestParameter(requestParameter); err != nil {
			return err
		}
		if err := saveRequestParameters(rp.SubParameters, requestParameter.ID, api_id); err != nil {
			return err
		}
	}

	return nil
}

func saveResponseParameters(rps []*ResponseParameters, p_id uint, api_id uint) error {
	if len(rps) == 0 {
		return nil
	}

	for _, rp := range rps {
		responseParameter := new(models.ResponseParameter)
		copier.Copy(responseParameter, rp)
		responseParameter.ID = 0
		responseParameter.ApiID = api_id
		responseParameter.ParentID = p_id
		fmt.Printf("-------%v\n", responseParameter)
		if err := models.CreateResponseParameter(responseParameter); err != nil {
			return err
		}
		if err := saveResponseParameters(rp.SubParameters, responseParameter.ID, api_id); err != nil {
			return err
		}
	}

	return nil
}

//func GetJson(c echo.Context) error {
//	apiJson := `{"id":1,"name":"用户","description":"","creator":0,"project":0,"group":0,"uri":"/users","Protocol":1,"Method":1,"request":{"headers":[{"id":1,"name":"Authorization","value":"Bearer eyJhbGciOiJI","description":"认证信息","api_id":1}],"parameters":[]},"response":{"headers":[{"id":2,"name":"Content-Type","value":"application/json","description":"数据响应的格式","remark":"","api_id":1}],"parameters":[{"id":1,"name":"users","value":"","type":7,"required":true,"description":"","remark":"用户集合","api_id":1,"parent_id":0,"sub_parameter":[{"id":2,"name":"username","value":"","type":2,"required":true,"description":"","remark":"用户名","api_id":1,"parent_id":1,"sub_parameter":null}]}]}}`
//	//return c.JSON(http.StatusOK, )

//	js, err := simplejson.NewJson([]byte(apiJson))
//	if err != nil {
//		fmt.Println("error")
//	}

//	responseJson := js.Get("response").Get("parameters")
//	if responseJson == nil {
//		fmt.Println("error")
//	}

//	parameters, _ := responseJson.Array()

//	for _, p := range parameters {
//		fmt.Println(p)
//	}

//	return c.JSON(http.StatusOK, parameters)
//}
