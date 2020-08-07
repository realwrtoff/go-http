package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type DistinctReq struct {
	Business string `form:"business"`
	Key      string `form:"key"`
}

func (s *Service) Distinct(c *gin.Context) {
	req := &DistinctReq{}
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

	key := fmt.Sprintf("%s:%s", req.Business, req.Key)
	curTime := time.Now().Format("2006-01-02")
	ok, err := s.rds.SetNX(key, curTime, time.Hour*24*365).Result()
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	} else if ok {
		res.Message = fmt.Sprintf("setnx %s success", key)
		res.Data = curTime
	} else {
		res.Message = fmt.Sprintf("setnx %s failed", key)
	}
	c.JSON(http.StatusOK, res)
}
