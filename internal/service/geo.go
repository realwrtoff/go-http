package service

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
)

type GeoPointReq struct {
	Latitude1  float64 `form:"latitude_1" json:"latitude_1"`
	Longitude1 float64 `form:"longitude_1" json:"longitude_1"`
	Latitude2  float64 `form:"latitude_2" json:"latitude_2"`
	Longitude2 float64 `form:"longitude_2" json:"longitude_2"`
}

type Point struct {
	Latitude  float64 `form:"latitude" json:"lat"`
	Longitude float64 `form:"longitude" json:"lng"`
}

type AddressReq struct {
	Key string `form:"key, omitempty" json:"key"`
	City string `form:"city, omitempty" json:"city"`
	Address string `form:"address" json:"address"`
}

type AdInfo struct {
	AdCode string `json:"adcode"`
}
type AddressComponent struct {
	Province string `json:"province"`
	City string `json:"city"`
	District string `json:"district"`
	Street string `json:"street"`
	StreetNumber string `json:"street_number"`
}

type AddressInfo struct {
	Title string `json:"title"`
	Location Point `json:"location"`
	AdInfo AdInfo `json:"ad_info"`
	AddressComponents AddressComponent `json:"address_components"`
	Reliability int32 `json:"reliability"`
}

type TxRes struct {
	Status int32 `json:"status"`
	Message string `json:"message"`
	Result AddressInfo `json:"result"`
}

func (s *Service) Distance(c *gin.Context) {
	req := &GeoPointReq{}
	res := &BasicRes{
		Code:    200,
		Message: "",
		Data:    nil,
	}

	if err := c.Bind(req); err != nil {
		res.Code = http.StatusBadRequest
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	radius := 6371000.0 //6378137.0
	rad := math.Pi / 180.0
	lat1 := req.Latitude1 * rad
	lng1 := req.Longitude1 * rad
	lat2 := req.Latitude2 * rad
	lng2 := req.Longitude2 * rad
	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	distance := dist * radius / 1000
	res.Data = distance
	c.JSON(http.StatusOK, res)
}

func (s *Service) Address(c *gin.Context) {
	req := &AddressReq{}
	res := &BasicRes{
		Code:    200,
		Message: "",
		Data:    nil,
	}

	if err := c.Bind(req); err != nil {
		res.Code = http.StatusBadRequest
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	var key string
	if req.Key == "" {
		key = "V5NBZ-DZOKO-WYFWD-SIKJJ-BDYY5-ZVF6I"
	} else {
		key = req.Key
	}
	txReq := make(map[string]interface{})
	txReq["key"] = key
	txReq["region"] = req.City
	txReq["address"] = req.Address
	uri := "https://apis.map.qq.com/ws/geocoder/v1"
	httpRes := s.httpClient.GET(uri, nil, txReq, nil)
	if httpRes.Err != nil {
		s.runLog.Errorf("request %s params %v, response [%d][%s]", uri, txReq, httpRes.Status, httpRes.Err.Error())
		res.Message = httpRes.Err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	txRes := TxRes{}
	s.runLog.Info(httpRes.Res)
	_ = json.Unmarshal(httpRes.Res, &txRes)
	res.Code = int(txRes.Status)
	res.Message = txRes.Message
	res.Data = txRes.Result
	c.JSON(http.StatusOK, res)
}
