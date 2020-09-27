package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func BKDRHash(str string) uint64 {
	seed := uint64(131) // 31 131 1313 13131 131313 etc..
	hash := uint64(0)
	for i := 0; i < len(str); i++ {
		hash = (hash * seed) + uint64(str[i])
	}
	return hash & 0x7FFFFFFF
}

type MobileSetReq struct {
	MobilePrefix string `form:"mobile_prefix"`
}

type MobileQueryReq struct {
	MobileMd5 string `form:"mobile_md5"`
}

type MobileQueryMultiReq struct {
	MobileMd5 []string `form:"mobile_md5"`
}

func Md5V(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func TimeCost(start time.Time) time.Duration{
	return time.Since(start)
}

func (s *Service) SetMobileMd5(c *gin.Context) {
	req := &MobileSetReq{}
	res := &BasicRes{
		Code:    http.StatusOK,
		Message: "",
		Data:    nil,
	}

	if err := c.Bind(req); err != nil {
		res.Code = http.StatusBadRequest
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	okNum := 0
	start := time.Now()
	for i := 0; i < 10000; i++ {
		mobile := fmt.Sprintf("%s%04d", req.MobilePrefix, i)
		mobileMd5 := Md5V(mobile)
		key := fmt.Sprintf("%d", s.crc32.Hash32S(mobileMd5) % 10000000)
		field := fmt.Sprintf("%d", BKDRHash(mobileMd5))
		ok, err := s.rds.HSetNX(key, field, mobile).Result()
		if err != nil {
			s.runLog.Errorf("hset %s %s %s error[%s]", key, field, mobile, err.Error())
			continue
		}
		if !ok {
			s.runLog.Errorf("hset %s %s %s failed", key, field, mobile)
			continue
		}
		okNum += 1
	}
	duration := TimeCost(start)
	res.Message = fmt.Sprintf("hsetnx mobile[%s] success %d, time spend %d milliseconds ", req.MobilePrefix, okNum, duration/1000000)
	s.runLog.Info(res.Message)
	c.JSON(http.StatusOK, res)
}

func (s *Service) QueryMobile(c *gin.Context) {
	req := &MobileQueryReq{}
	res := &BasicRes{
		Code:    http.StatusOK,
		Message: "success",
		Data:    nil,
	}

	if err := c.Bind(req); err != nil {
		res.Code = http.StatusBadRequest
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	key := fmt.Sprintf("%d", s.crc32.Hash32S(req.MobileMd5) % 10000000)
	field := fmt.Sprintf("%d", BKDRHash(req.MobileMd5))
	mobile, err := s.rds.HGet(key, field).Result()
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
	} else {
		res.Data = mobile
	}
	c.JSON(http.StatusOK, res)
}

func (s *Service) QueryMobileMulti(c *gin.Context) {
	req := &MobileQueryMultiReq{}
	res := &BasicRes{
		Code:    http.StatusOK,
		Message: "success",
		Data:    nil,
	}

	if err := c.Bind(req); err != nil {
		res.Code = http.StatusBadRequest
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}
	data := make(map[string]string)
	for _, mobileMd5 := range req.MobileMd5 {
		key := fmt.Sprintf("%d", s.crc32.Hash32S(mobileMd5) % 10000000)
		field := fmt.Sprintf("%d", BKDRHash(mobileMd5))
		mobile, err := s.rds.HGet(key, field).Result()
		if err != nil {
			s.runLog.Errorf("md5 %s hget %s %s error[%s]", mobileMd5, key, field, err.Error())
		} else {
			data[mobileMd5] = mobile
		}
	}
	res.Data = data
	c.JSON(http.StatusOK, res)
}