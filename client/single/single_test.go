package client

import (
	"fmt"
	"go-Redisson/options"
	"testing"
)

func TestClient(t *testing.T) {
	client := NewClient(&options.Options{
		Addr:     "127.0.0.1:6379",
		Password: "123456",
		DB:       0})
	s, err := client.Set("test", "test")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(s)
	//c := &config.Config{
	//	Addr:     "127.0.0.1:6379",
	//	Password: "123456",
	//	DB:       0,
	//}
	//client, err := NewClient(c)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//s, err := client.Get("1")
	//fmt.Println(s, err)
	//s, err = client.Set("2", "2")
	//fmt.Println(s, err)
	//s, err = client.Del("2")
	//fmt.Println(s, err)
	//s, err = client.Get("2")
	//fmt.Println(s, err)
	//c.DB = 1
	//client2, err := NewClient(c)
	//if err != nil {
	//	fmt.Println(err)
	//}
	//s, err = client2.Set("2", "2")
	//fmt.Println(s, err)

	//client.Conn.Close()
}
