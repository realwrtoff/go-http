package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/hpifu/go-kit/hashx"
	"net/http"
	"time"
)


const ConflictMd5 = "conflict_md5"
const RedisKeyBudgets = 10000000

type MobileSetReq struct {
	MobilePrefix string `form:"mobile_prefix"`
}

type MobileQueryReq struct {
	MobileMd5 string `form:"mobile_md5"`
}

type MobileQueryMultiReq struct {
	MobileMd5 []string `form:"mobile_md5"`
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
		mobileMd5 := hashx.Md5Hash(mobile)
		key := fmt.Sprintf("%d", s.crc.Hash64S(mobileMd5) % RedisKeyBudgets)
		field := fmt.Sprintf("%d", hashx.DJBHash(mobileMd5))
		ok, err := s.rds.HSetNX(key, field, mobile).Result()
		if err != nil {
			s.runLog.Errorf("hset %s %s %s error[%s]", key, field, mobile, err.Error())
			continue
		}
		if !ok {
			s.runLog.Errorf("hset %s %s %s failed", key, field, mobile)
			// 放入一个固定key的hash中
			s.rds.HSet(ConflictMd5, mobileMd5, mobile)
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

	var mobile string
	mobile, err := s.rds.HGet(ConflictMd5, req.MobileMd5).Result()
	if err == redis.Nil {
		key := fmt.Sprintf("%d", s.crc.Hash64S(req.MobileMd5) % RedisKeyBudgets)
		field := fmt.Sprintf("%d", hashx.DJBHash(req.MobileMd5))
		mobile, err = s.rds.HGet(key, field).Result()
	}
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
		var mobile string
		mobile, err := s.rds.HGet(ConflictMd5, mobileMd5).Result()
		if err == redis.Nil {
			key := fmt.Sprintf("%d", s.crc.Hash64S(mobileMd5) % RedisKeyBudgets)
			field := fmt.Sprintf("%d", hashx.DJBHash(mobileMd5))
			mobile, err = s.rds.HGet(key, field).Result()
		}
		if err == nil {
			data[mobileMd5] = mobile
		}
	}
	res.Data = data
	c.JSON(http.StatusOK, res)
}