package service

import (
	"github.com/go-redis/redis"
	"github.com/hpifu/go-kit/hhttp"
	"github.com/sirupsen/logrus"
)

type Service struct {
	rds    *redis.Client
	httpClient *hhttp.HttpClient
	runLog *logrus.Logger
}

func NewService(
	rds *redis.Client,
	httpClient *hhttp.HttpClient,
	runLog *logrus.Logger,
) *Service {
	return &Service{
		rds:    rds,
		httpClient: httpClient,
		runLog: runLog,
	}
}

type BasicRes struct {
	Code    int         `form:"code" json:"code"`
	Message string      `form:"message" json:"message"`
	Data    interface{} `form:"data" json:"data"`
}

type PageDataRes struct {
	Page       int         `json:"page"`
	TotalPage  int         `json:"totalPage"`
	TotalCount int         `json:"totalCount"`
	Items      interface{} `json:"items"`
}
