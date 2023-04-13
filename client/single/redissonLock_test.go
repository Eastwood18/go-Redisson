package client

import (
	"fmt"
	"go-Redisson/options"
	"testing"
)

func TestLockSet(t *testing.T) {
	client := NewClient(&options.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0})
	lock := NewLock()
	lock.lockName = "lockTest"
	err := client.Lock(lock)
	if err != nil {
		fmt.Println(err)
	}
}
