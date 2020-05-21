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
}

func (s *Service) Distinct(c *gin.Context) {
	req := &DistinctReq{}
	res := &BasicRes{
		Code: http.StatusOK,
		Message: "",
		Data: nil,
	}

	if err := c.Bind(req); err != nil {
		res.Code = http.StatusBadRequest
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	key := fmt.Sprintf("%s:%s", req.Business, req.Key)
	curTime := time.Now().Format("20060102150405")
	value, err := s.rds.GetSet(key, curTime).Result()
	if err == redis.Nil {
		res.Message = fmt.Sprintf("set %s first success", key)
	} else if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	} else {
		res.Message = fmt.Sprintf("get %s old value %s", key, value)
		res.Data = value
	}
	c.JSON(http.StatusOK, res)
}
