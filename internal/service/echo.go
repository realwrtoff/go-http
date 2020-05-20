package service

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type EchoReq struct {
	Message string `form:"message"`
}

type EchoRes struct {
	Message string `form:"message" json:"message"`
}

type MobileReq struct {
	Mobile string `form:"mobile"`
}

type MobileInfo struct {
	Province string `form:"province" json:"province"`
	City     string `form:"city" json:"city"`
	Operator string `form:"operator" json:"operator"`
	PostCode string `form:"post_code" json:"post_code"`
	AreaCode string `form:"area_code" json:"area_code"`
}

type MobileRes struct {
	Message string `form:"message" json:"message"`
	Data *MobileInfo `form:"data" json:"data"`
}

func (s *Service) Echo(c *gin.Context) {
	req := &EchoReq{}

	if err := c.Bind(req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	res := &EchoRes{
		Message: req.Message,
	}
	c.JSON(http.StatusOK, res)
}

func (s *Service) Location(c *gin.Context) {
	req := &MobileReq{}
	res := &MobileRes{
		Message: "",
		Data: nil,
	}

	if err := c.Bind(req); err != nil {
		res.Message = err.Error()
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if len(req.Mobile) < 7 {
		res.Message = fmt.Sprintf("invalid mobile prefix [%s], at least length 7.", req.Mobile)
		c.JSON(http.StatusBadRequest, res)
		return
	}

	key := fmt.Sprintf("%s", req.Mobile[:7])
	value, err := s.rds.Get(key).Result()
	if err != nil {
		res.Message = err.Error()
		c.JSON(http.StatusInternalServerError, res)
		return
	}
	items := strings.Split(value, "|")
	data := &MobileInfo{}
	data.Province = items[0]
	data.City = items[1]
	data.PostCode = items[2]
	data.AreaCode = items[3]
	data.Operator = items[4]

	// value, err := s.rds.Get(key).Bytes()
	//if err := json.Unmarshal(value, data); err != nil {
	//	res.Message = err.Error()
	//	c.JSON(http.StatusInternalServerError, res)
	//}

	res.Data = data
	c.JSON(http.StatusOK, res)
}