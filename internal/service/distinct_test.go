package service

import (
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestService_Distinct(t *testing.T) {
	rds := redis.NewClient(&redis.Options{Addr: "52.130.89.55:21603", Password: "Rds123"})
	value, err := rds.SetNX("jim", 1, 10*time.Second).Result()
	if err != nil {
		t.Errorf("setnx err %v", err.Error())
	} else {
		t.Logf("setnx return %v", value)
	}
	value, err = rds.SetNX("jim", 1, 10*time.Second).Result()
	if err != nil {
		t.Errorf("setnx again err %v", err.Error())
	} else {
		t.Logf("setnx again return %v", value)
	}
}
