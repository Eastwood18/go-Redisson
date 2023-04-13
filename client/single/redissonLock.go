package client

import (
	"context"
	"fmt"
	"github.com/gofrs/uuid/v5"
	"time"
)

type RedisLock struct {
	lockName  string
	createdAt time.Time
	duration  time.Duration
	count     int
	value     string
	thread    string
}

func NewLock() *RedisLock {
	l := new(RedisLock)
	l.createdAt = time.Now()
	thread, _ := uuid.NewV4()
	l.thread = thread.String()
	l.count = 0
	return l
}
func (c *Client) Lock(redisLock *RedisLock) (err error) {
	ctx := context.Background()
	conn, err := c.getConn(ctx)
	do, err := conn.Do(ctx, "exists", redisLock.lockName)
	if err != nil {
		return
	}
	fmt.Println(do)
	if do == "0" {
		_, err = conn.Do(ctx, "hmset", redisLock.lockName, "thread", redisLock.thread, "count", "1")
		if err != nil {
			return
		}
	} else if do == "1" {
		_, err = conn.Do(ctx, "HINCRBY", redisLock.lockName, "count", "1")
		if err != nil {
			return
		}
	}
	return
}

//func (c *Client) GetLock(redisLock *RedisLock) (err error) {
//	ticker := time.NewTicker(redisLock.duration)
//
//	select {
//	case <-ticker.C:
//
//	}
//}
