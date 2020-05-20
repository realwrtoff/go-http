package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"time"
)

type DistinctReq struct {
	Business string `form:"business"`
	Key string `form:"key"`
	Operation string `form:"operation"`
}

func (s *Service) Distinct(c *gin.Context) {
	req := &DistinctReq{}
	res := &BasicRes{
		Code: 0,
		Message: "",
		Data: nil,
	}

	if err := c.Bind(req); err != nil {
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	key := fmt.Sprintf("%s:%s", req.Business, req.Key)
	if req.Operation == "set" {
		curTime := time.Now().Format("20060102150405")
		if _, err := s.rds.Set(key, curTime, time.Hour*24*365).Result(); err != nil {
			res.Code = http.StatusInternalServerError
			res.Message = err.Error()
			c.JSON(http.StatusInternalServerError, res)
			return
		} else {
			res.Message = fmt.Sprintf("set %s datetime %s success", key, curTime)
		}
	} else if req.Operation == "get" {
		value, err := s.rds.Get(key).Result()
		if err == redis.Nil {
			res.Code = http.StatusNotFound
			res.Message = fmt.Sprintf("%s not found", key)
		} else if err != nil {
			res.Code = http.StatusInternalServerError
			res.Message = err.Error()
			c.JSON(http.StatusInternalServerError, res)
			return
		} else {
			res.Message = fmt.Sprintf("%s found, datetime %s", key, value)
		}
	} else {
		curTime := time.Now().Format("20060102150405")
		ok, err := s.rds.SetNX(key, curTime, time.Hour*24*365).Result()
		if err != nil {
			res.Code = http.StatusInternalServerError
			res.Message = err.Error()
			c.JSON(http.StatusInternalServerError, res)
			return
		}
		if ok {
			res.Message = fmt.Sprintf("setnx %s datetime %s success", key, curTime)
		} else {
			res.Code = -1
			res.Message = fmt.Sprintf("setnx %s failed", key)
			c.JSON(http.StatusAlreadyReported, res)
			return
		}
	}
	c.JSON(http.StatusOK, res)
}
