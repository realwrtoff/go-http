package monitor

import (
	"github.com/go-redis/redis"
	"time"
)

type Watcher struct {
	rds *redis.Client
}


func NewWatcher(rds *redis.Client) *Watcher {
	return &Watcher{
		rds:rds,
	}
}

func (w *Watcher) Run() {
	go func() {
		ticker := time.Tick(time.Minute)
		for {
			select {
			case <-ticker:
				w.rds.Ping()
			}
		}
	}()
}
