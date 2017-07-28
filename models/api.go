package models

import (
	"fmt"
	//"errors"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	log "github.com/sirupsen/logrus"
)

const (
	// parameter type
	_ = iota
	Number
	String
	Object
	Boolean
	ArrayNumber
	ArrayString
	ArrayObject
	ArrayBoolean
	Array
)

const (
	// request method
	_ = iota
	GET
	POST
	PUT
	DELETE
	HEAD
	PATCH
	OPTIONS
)

const (
	// request protocol
	_ = iota
	HTTP
	HTTPS
)

type Api struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Name        string    `json:"name" gorm:"not null"`
	Description string    `json:"description" validate:"max=100"`
	CreatorID   uint      `json:"creator" gorm:"default:0"`
	ProjectID   uint      `json:"project" gorm:"default:0"`
	GroupID     uint      `json:"group" gorm:"default:0"`
	URI         string    `json:"uri" gorm:"not null" validate:"required,max=100"`
	Protocol    uint      `json:"protocol" validate:"required,max=20"`
	Method      uint      `json:"method" validate:"required,max=20"`
}

type RequestHeader struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Name        string    `json:"name" gorm:"not null" validate:"required,max=20"`
	Value       string    `json:"value" gorm:"not null"`
	Required    bool      `json:"required" gorm:"default:0"`
	Description string    `json:"description" validate:"max=100"`
	ApiID       uint      `json:"api_id" gorm:"not null"`
}

type ResponseHeader struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Name        string    `json:"name" gorm:"not null" validate:"required,max=20"`
	Value       string    `json:"value" gorm:"not null"`
	Required    bool      `json:"required" gorm:"default:0"`
	Description string    `json:"description" validate:"max=100"`
	Remark      string    `json:"remark"`
	ApiID       uint      `json:"api_id" gorm:"not null"`
}

type RequestParameter struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Name        string    `json:"name" gorm:"not null" validate:"required,max=20"`
	Value       string    `json:"value" gorm:"not null"`
	Type        uint      `json:"type"`
	Required    bool      `json:"required" gorm:"default:0"`
	Description string    `json:"description" validate:"max=100"`
	Remark      string    `json:"remark" validate:"required,max=100"`
	ApiID       uint      `json:"api_id" gorm:"not null"`
	ParentID    uint      `json:"parent_id"`
}

type ResponseParameter struct {
	ID          uint      `json:"id" gorm:"primary_key"`
	CreatedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	Name        string    `json:"name" gorm:"not null" validate:"required,max=20"`
	Value       string    `json:"value" gorm:"not null"`
	Type        uint      `json:"type"`
	Required    bool      `json:"required" gorm:"default:0"`
	Description string    `json:"description" validate:"max=100"`
	Remark      string    `json:"remark" validate:"required,max=100"`
	ApiID       uint      `json:"api_id" gorm:"not null"`
	ParentID    uint      `json:"parent_id"`
}

func (api *Api) AfterSave() (err error) {
	d := new(ApiIndex)
	copier.Copy(d, api)
	d.SearchType = "api"
	d.ID = strconv.Itoa(int(api.ID))
	err = BleveIndex.Index("user:"+d.ID, d)
	return
}

func (api *Api) AfterDelete() (err error) {
	err = BleveIndex.Delete("api:" + strconv.Itoa(int(api.ID)))
	return
}

// api
func CreateOrUpdateApi(api *Api) error {
	//err := db.Create(api).Error
	err := db.Save(api).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":  err.Error(),
			"api": *api,
		}).Error("create api error")
		return err
	}

	if err := DeleteApiRequestInfoByID(api.ID); err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"api": *api,
	}).Info("create api success")

	return nil
}

func GetApiByID(id uint) (*Api, error) {
	api := new(Api)
	err := db.First(api, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("get api error")
		return nil, err
	}

	return api, nil
}

func UpdateApi(api *Api) error {
	err := db.Model(api).Updates(api).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":  err.Error(),
			"api": *api,
		}).Error("update api error")
		return err
	}

	log.WithFields(log.Fields{
		"api": *api,
	}).Info("update api group success")

	return nil
}

func DeleteApiByID(id uint) error {
	err := db.Where("id = ?", id).Delete(Api{}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("delete api error")
		return err
	}

	if err := DeleteApiRequestInfoByID(id); err != nil {
		return err
	}

	return nil
}

func DeleteApiRequestInfoByID(id uint) error {
	tx := db.Begin()
	tx.Where("api_id = ?", id).Delete(RequestHeader{})
	tx.Where("api_id = ?", id).Delete(RequestParameter{})
	tx.Where("api_id = ?", id).Delete(ResponseHeader{})
	tx.Where("api_id = ?", id).Delete(ResponseParameter{})
	if err := tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
		}).Error("delete api request info fail")
		return err
	}

	log.Info("delete api request info success")

	return nil
}

// request header
func BatchCreateRequestHeader(rhs []*RequestHeader) error {
	tx := db.Begin()
	for _, rh := range rhs {
		tx.Create(rh)
	}
	if err := tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
		}).Error("batch create requestHeader error")
		return err
	}

	log.Info("batch create requestHeader success")

	return nil
}

func CreateRequestHeader(requestHeader *RequestHeader) error {
	err := db.Create(requestHeader).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":            err.Error(),
			"requestHeader": *requestHeader,
		}).Error("create requestHeader error")
		return err
	}

	log.WithFields(log.Fields{
		"requestHeader": *requestHeader,
	}).Info("create requestHeader success")

	return nil
}

func GetApiRequestHeadersByID(id uint) ([]*RequestHeader, error) {
	rhs := make([]*RequestHeader, 0)
	err := db.Where("api_id = ?", id).Find(&rhs).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return rhs, nil
}

func GetRequestHeaderByID(id uint) (*RequestHeader, error) {
	requestHeader := new(RequestHeader)
	err := db.First(requestHeader, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("get requestHeader error")
		return nil, err
	}

	return requestHeader, nil
}

func UpdateRequestHeader(requestHeader *RequestHeader) error {
	err := db.Model(requestHeader).Updates(requestHeader).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":            err.Error(),
			"requestHeader": *requestHeader,
		}).Error("update requestHeader error")
		return err
	}

	log.WithFields(log.Fields{
		"requestHeader": *requestHeader,
	}).Info("update requestHeader group success")

	return nil
}

func DeleteRequestHeaderByID(id uint) error {
	err := db.Where("id = ?", id).Delete(RequestHeader{}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("delete requestHeader error")
		return err
	}

	return nil
}

// response header
func BatchCreateResponseHeader(rhs []*ResponseHeader) error {
	tx := db.Begin()
	for _, rh := range rhs {
		tx.Create(rh)
	}
	if err := tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
		}).Error("batch create responseHeader error")
		return err
	}

	log.Info("batch create responseHeader success")

	return nil
}

func CreateResponseHeader(responseHeader *ResponseHeader) error {
	err := db.Create(responseHeader).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":             err.Error(),
			"responseHeader": *responseHeader,
		}).Error("create responseHeader error")
		return err
	}

	log.WithFields(log.Fields{
		"responseHeader": *responseHeader,
	}).Info("create responseHeader success")

	return nil
}

func GetApiResponseHeadersByID(id uint) ([]*ResponseHeader, error) {
	rhs := make([]*ResponseHeader, 0)
	err := db.Where("api_id = ?", id).Find(&rhs).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return rhs, nil
}

func GetResponseHeaderByID(id uint) (*ResponseHeader, error) {
	responseHeader := new(ResponseHeader)
	err := db.First(responseHeader, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("get responseHeader error")
		return nil, err
	}

	return responseHeader, nil
}

func UpdateResponseHeader(responseHeader *ResponseHeader) error {
	err := db.Model(responseHeader).Updates(responseHeader).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":             err.Error(),
			"responseHeader": *responseHeader,
		}).Error("update responseHeader error")
		return err
	}

	log.WithFields(log.Fields{
		"responseHeader": *responseHeader,
	}).Info("update responseHeader group success")

	return nil
}

func DeleteResponseHeaderByID(id uint) error {
	err := db.Where("id = ?", id).Delete(ResponseHeader{}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("delete responseHeader error")
		return err
	}

	return nil
}

// request parameter
func BatchCreateRequestParameter(rps []*RequestParameter) error {
	tx := db.Begin()
	for _, rp := range rps {
		tx.Create(rp)
	}
	if err := tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
		}).Error("batch create request parameter error")
		return err
	}

	log.Info("batch create request parameter success")

	return nil
}

func CreateRequestParameter(requestParameter *RequestParameter) error {
	err := db.Create(requestParameter).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":               err.Error(),
			"requestParameter": *requestParameter,
		}).Error("create request parameter error")
		return err
	}

	log.WithFields(log.Fields{
		"requestParameter": *requestParameter,
	}).Info("create request parameter success")

	return nil
}

func GetRequestHeadersByID(api_id, parent_id uint) ([]*RequestParameter, error) {
	rps := make([]*RequestParameter, 0)
	err := db.Where("api_id = ? and parent_id = ?", api_id, parent_id).Find(&rps).Error
	if err != nil {
		return nil, err
	}

	return rps, nil
}

func GetRequestParameterByID(id uint) (*RequestParameter, error) {
	requestParameter := new(RequestParameter)
	err := db.First(requestParameter, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("get request parameter error")
		return nil, err
	}

	return requestParameter, nil
}

func UpdateRequestParameter(requestParameter *RequestParameter) error {
	err := db.Model(requestParameter).Updates(requestParameter).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":               err.Error(),
			"requestParameter": *requestParameter,
		}).Error("update request parameter error")
		return err
	}

	log.WithFields(log.Fields{
		"request parameter": *requestParameter,
	}).Info("update request parameter group success")

	return nil
}

func DeleteRequestParameterByID(id uint) error {
	err := db.Where("id = ?", id).Delete(RequestParameter{}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("delete request parameter error")
		return err
	}

	return nil
}

// response parameter
func BatchCreateResponseParameter(rps []*ResponseParameter) error {
	tx := db.Begin()
	for _, rp := range rps {
		tx.Create(rp)
	}
	if err := tx.Commit().Error; err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
		}).Error("batch create response parameter error")
		return err
	}

	log.Info("batch create response parameter success")

	return nil
}

func CreateResponseParameter(responseParameter *ResponseParameter) error {
	err := db.Create(responseParameter).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":                err.Error(),
			"responseParameter": *responseParameter,
		}).Error("create response parameter error")
		return err
	}

	log.WithFields(log.Fields{
		"responseParameter": *responseParameter,
	}).Info("create response parameter success")

	return nil
}

func GetResponseHeadersByID(api_id, parent_id uint) ([]*ResponseParameter, error) {
	rps := make([]*ResponseParameter, 0)
	err := db.Where("api_id = ? and parent_id = ?", api_id, parent_id).Find(&rps).Error
	if err != nil {
		return nil, err
	}

	return rps, nil
}

func GetResponseParameterByID(id uint) (*ResponseParameter, error) {
	responseParameter := new(ResponseParameter)
	err := db.First(responseParameter, id).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("get response parameter error")
		return nil, err
	}

	return responseParameter, nil
}

func UpdateResponseParameter(responseParameter *ResponseParameter) error {
	err := db.Model(responseParameter).Updates(responseParameter).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db":                err.Error(),
			"responseParameter": *responseParameter,
		}).Error("update response parameter error")
		return err
	}

	log.WithFields(log.Fields{
		"response parameter": *responseParameter,
	}).Info("update response parameter group success")

	return nil
}

func DeleteResponseParameterByID(id uint) error {
	err := db.Where("id = ?", id).Delete(ResponseParameter{}).Error
	if err != nil {
		log.WithFields(log.Fields{
			"db": err.Error(),
			"id": id,
		}).Error("delete response parameter error")
		return err
	}

	return nil
}
