package service

import (
	"github.com/go-redis/redis"
	"github.com/hpifu/go-kit/hashx"
	"github.com/hpifu/go-kit/hhttp"
	"github.com/sirupsen/logrus"
)

type Service struct {
	rds    *redis.Client
	httpClient *hhttp.HttpClient
	crc hashx.Hasher64
	runLog *logrus.Logger
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

func (s *Service)Init() error{
	s.crc = hashx.NewHasher64(hashx.CRC64ISO)
	return nil
}