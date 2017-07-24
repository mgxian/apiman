package api

import (
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
	//"github.com/bitly/go-simplejson"
	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
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

	t, _ := models.GetTeamByID(p.Team)
	if t == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "team不存在",
		})
	}

	//username := tokenInfo.Name
	teamname := t.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag {
		flag = models.IsTeamMember(teamname, username)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	userid := strconv.Itoa(int(u.ID))
	apif.Creator = userid
	apif.Group = g.ID
	apif.Project = p.ID

	err = saveApi(apif, true)
	if err != nil {
		log.WithFields(log.Fields{
			"operator": tokenInfo.Name,
			"error":    err.Error(),
		}).Info("create or update api fail")
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"operator": tokenInfo.Name,
		"api":      *apif,
	}).Info("create or update api success")

	return c.JSON(http.StatusCreated, apif)
}

func GetApi(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	api_idstr := c.Param("id")
	api_id, _ := strconv.Atoi(api_idstr)

	api := new(ApiForm)

	// get base info
	apiBaseInfo, _ := models.GetApiByID(uint(api_id))
	if apiBaseInfo == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "API不存在",
		})
	}

	username := tokenInfo.Name
	u, _ := models.GetUserByName(username)
	if u == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "用户不存在",
		})
	}

	p, _ := models.GetProjectByID(apiBaseInfo.Project)
	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "项目不存在",
		})
	}

	t, _ := models.GetTeamByID(p.Team)
	if t == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "team不存在",
		})
	}

	//username := tokenInfo.Name
	teamname := t.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag {
		flag = models.IsTeamMember(teamname, username)
	}

	if !flag {
		flag = models.IsTeamReader(teamname, username)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	copier.Copy(api, apiBaseInfo)
	ac, _ := models.GetUserByID(apiBaseInfo.Creator)
	api.Creator = ac.Name

	// get header
	requestHeaders, _ := models.GetApiRequestHeadersByID(uint(api_id))
	responseHeaders, _ := models.GetApiResponseHeadersByID(uint(api_id))
	api.Request.RequestHeaders = requestHeaders
	api.Response.ResponseHeaders = responseHeaders

	//fmt.Printf("----%v", requestHeaders)
	//fmt.Println(tokenInfo.Name)

	// get parameter
	req := getRequestParameters(api.ID, uint(0))
	resp := getResponseParameters(api.ID, uint(0))
	api.Request.RequestParameters = req
	api.Response.ResponseParameters = resp

	return c.JSONPretty(http.StatusOK, api, "  ")
}

func UpdateApi(c echo.Context) error {
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

	api_id := c.Param("id")
	api_idstr, _ := strconv.Atoi(api_id)
	apif.ID = uint(api_idstr)

	api, _ := models.GetApiByID(uint(api_idstr))
	if api == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "api不存在",
		})
	}

	g, _ := models.GetApiGroupByID(api.Group)
	if g == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "api group不存在",
		})
	}

	username := tokenInfo.Name

	p, _ := models.GetProjectByID(g.Project)
	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "项目不存在",
		})
	}

	t, _ := models.GetTeamByID(p.Team)
	if t == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "team不存在",
		})
	}

	//username := tokenInfo.Name
	teamname := t.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag {
		flag = models.IsTeamMember(teamname, username)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	apif.ID = api.ID
	apif.CreatedAt = api.CreatedAt
	apif.Creator = strconv.Itoa(int(api.Creator))
	apif.Group = api.Group
	apif.Project = api.Project

	err = saveApi(apif, false)

	if err != nil {
		log.WithFields(log.Fields{
			"operator": tokenInfo.Name,
			"error":    err.Error(),
		}).Info("update api fail")
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"operator": tokenInfo.Name,
		"api":      *apif,
	}).Info("update api success")

	return c.JSON(http.StatusCreated, apif)
}

func DeleteApi(c echo.Context) error {
	tokenInfo, err := jwt.GetClaims(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": err.Error(),
			})
	}

	//fmt.Println(tokenInfo.Name)

	id := c.Param("id")
	intstr, _ := strconv.Atoi(id)

	api, _ := models.GetApiByID(uint(intstr))
	if api == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "api不存在",
		})
	}

	p, _ := models.GetProjectByID(api.Project)
	if p == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "project不存在",
		})
	}

	t, _ := models.GetTeamByID(p.Team)
	if t == nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "team不存在",
		})
	}

	username := tokenInfo.Name
	teamname := t.Name

	flag := models.IsTeamMaintainer(teamname, username)

	if !flag {
		flag = models.IsTeamMember(teamname, username)
	}

	if !flag && !tokenInfo.Admin {
		return c.JSON(http.StatusUnauthorized,
			echo.Map{
				"message": "你没有此权限",
			})
	}

	if err := models.DeleteApiByID(uint(intstr)); err != nil {
		log.WithFields(log.Fields{
			"operator": username,
			"error":    err.Error(),
			"api":      *api,
		}).Error("delete api fail")
		return c.NoContent(http.StatusInternalServerError)
	}

	log.WithFields(log.Fields{
		"api":      *api,
		"operator": username,
	}).Error("delete api success")
	return c.NoContent(http.StatusNoContent)
}

func getRequestParameters(api_id, p_id uint) []*RequestParameters {
	rps_m, _ := models.GetRequestHeadersByID(api_id, p_id)
	if len(rps_m) == 0 {
		//fmt.Println("-----------rps_m nil-----------------")
		return make([]*RequestParameters, 0)
	}

	//fmt.Printf("%v", rps_m)

	rps := make([]*RequestParameters, 0)
	for _, rp := range rps_m {
		t := new(RequestParameters)
		copier.Copy(t, rp)
		rps = append(rps, t)
	}

	if len(rps) == 0 {
		//fmt.Println("-----------rps nil-----------------")
		return make([]*RequestParameters, 0)
	}

	for _, rp := range rps {
		//fmt.Printf("%v", rp)
		rps_t := getRequestParameters(api_id, rp.ID)
		rp.SubParameters = rps_t
	}

	return rps
}

func getResponseParameters(api_id, p_id uint) []*ResponseParameters {
	rps_m, _ := models.GetResponseHeadersByID(api_id, p_id)
	if len(rps_m) == 0 {
		//fmt.Println("-----------rps_m nil-----------------")
		return make([]*ResponseParameters, 0)
	}

	//fmt.Printf("%v", rps_m)

	rps := make([]*ResponseParameters, 0)
	for _, rp := range rps_m {
		t := new(ResponseParameters)
		copier.Copy(t, rp)
		rps = append(rps, t)
	}

	if len(rps) == 0 {
		//fmt.Println("-----------rps nil-----------------")
		return make([]*ResponseParameters, 0)
	}

	for _, rp := range rps {
		//fmt.Printf("%v", rp)
		rps_t := getResponseParameters(api_id, rp.ID)
		rp.SubParameters = rps_t
	}

	return rps
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

	if err := models.CreateOrUpdateApi(apiBaseInfo); err != nil {
		return err
	}

	// save header
	for _, rh := range api.Request.RequestHeaders {
		//copier.Copy(apiRequestHeader, rh)
		rh.ID = 0
		rh.ApiID = apiBaseInfo.ID

		//fmt.Printf("%v", rh)
	}

	if err := models.BatchCreateRequestHeader(api.Request.RequestHeaders); err != nil {
		fmt.Print(err)
		return err
	}

	for _, rh := range api.Response.ResponseHeaders {
		//copier.Copy(apiResponseHeader, rh)
		rh.ID = 0
		rh.ApiID = apiBaseInfo.ID

		//fmt.Printf("%v", rh)
	}

	if err := models.BatchCreateResponseHeader(api.Response.ResponseHeaders); err != nil {
		fmt.Print(err)
		return err
	}

	// save parameter
	if err := saveRequestParameters(api.Request.RequestParameters, uint(0), apiBaseInfo.ID); err != nil {
		fmt.Print(err)
		return err
	}

	if err := saveResponseParameters(api.Response.ResponseParameters, uint(0), apiBaseInfo.ID); err != nil {
		fmt.Print(err)
		return err
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
		//fmt.Printf("-------%v\n", requestParameter)
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
		//fmt.Printf("-------%v\n", responseParameter)
		if err := models.CreateResponseParameter(responseParameter); err != nil {
			return err
		}
		if err := saveResponseParameters(rp.SubParameters, responseParameter.ID, api_id); err != nil {
			return err
		}
	}

	return nil
}
