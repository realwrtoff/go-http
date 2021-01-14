package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/realwrtoff/go-http/internal/public"
	"net/http"
	"strconv"
)

type BitMapReq struct {
	Business string `form:"business"`
	Key      string `form:"key"`
	Action   string `form:"action"`
}

func (s *Service) BitMap(c *gin.Context) {
	req := &BitMapReq{}
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
	if !public.IsValidMobile(req.Key) {
		res.Code = http.StatusBadRequest
		res.Message = fmt.Sprintf("mobile %s is invalid", req.Key)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	key := fmt.Sprintf("%s:%s", req.Business, req.Key[:5])
	offset, _ := strconv.ParseInt(req.Key[5:], 10, 64)
	var cnt int64
	var err error
	if req.Action == "set" {
		cnt, err = s.rds.SetBit(key, offset, 1).Result()
	} else {
		cnt, err = s.rds.GetBit(key, offset).Result()
	}
	if err != nil {
		res.Code = http.StatusInternalServerError
		res.Message = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	} else {
		res.Message = "success"
		res.Data = cnt
	}
	c.JSON(http.StatusOK, res)
}
