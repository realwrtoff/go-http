package service

import (
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

type Service struct {
	rds    *redis.Client
	runLog *logrus.Logger
}

func NewService(
	rds *redis.Client,
	runLog *logrus.Logger,
) *Service {
	return &Service{
		rds:    rds,
		runLog: runLog,
	}
}

type BasicRes struct {
	Code int `form:"code" json:"code"`
	Message string `form:"message" json:"message"`
	Data interface {} `form:"data" json:"data"`
}