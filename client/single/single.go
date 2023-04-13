package client

import (
	"context"
	"errors"
	"fmt"
	"go-Redisson/client/pool"
	"go-Redisson/options"
	"log"
	"net"
)

const (
	auth     = "auth%s\r\n"
	set      = "set%s\r\n"
	get      = "get%s\r\n"
	del      = "del%s\r\n"
	selectDB = "select%s\r\n"
)

type Client struct {
	opt      *options.Options
	connPool pool.Pooler
}

func NewClient(opt *options.Options) *Client {

	opt.Init()

	c := Client{
		opt: opt,
	}
	c.connPool = options.NewConnPool(opt, func(ctx context.Context, network, addr string) (net.Conn, error) {
		conn, err := net.Dial("tcp", opt.Addr)
		return conn, err
	})

	return &c
}

func (c *Client) newConn(ctx context.Context) (*pool.Conn, error) {
	cn, err := c.connPool.NewConn(ctx)
	if err != nil {
		return nil, err
	}

	err = c.initConn(ctx, cn)
	if err != nil {
		_ = c.connPool.CloseConn(cn)
		return nil, err
	}

	return cn, nil
}

func (c *Client) Set(key, value string) (result string, err error) {
	ctx := context.Background()
	conn, err := c.getConn(ctx)
	if err != nil {
		return "", err
	}
	conn.Write([]byte(fmt.Sprintf("set %s %s\r\n", key, value)))
	reader := conn.Reader()
	//defer reader.Reset()

	line, _, err := reader.ReadLine()
	if err != nil {
		log.Fatalln(err)
	}
	result = string(line[1:])
	return
}

//func (c *Client) Get(key string) (result string, err error) {
//	result, err = c.doCommand(get, key)
//	return
//}
//func (c *Client) Del(key string) (result string, err error) {
//	result, err = c.doCommand(del, key)
//	return
//}
//func (c *Client) doCommand(command string, options ...string) (result string, err error) {
//	ctx := context.Background()
//	opt := strings.Builder{}
//	for _, option := range options {
//		opt.WriteByte(' ')
//		opt.WriteString(option)
//	}
//	conn, err := c.connPool.Get(ctx)
//	if err != nil {
//		return "", err
//	}
//	_, err = c.Conn.Write([]byte(fmt.Sprintf(command, opt.String())))
//	if err != nil {
//		return "", err
//	}
//	reader := bufio.NewReader(c.Conn)
//	defer reader.Reset(c.Conn)
//	line, _, err := reader.ReadLine()
//	if err != nil {
//		log.Fatalln(err)
//	}
//	switch string(line[0]) {
//	case "+":
//		err = nil
//	case "-":
//		result, err = "", errors.New(string(line[1:]))
//		return
//	case ":":
//		err = nil
//	case "$":
//		err = nil
//	case "*":
//		err = nil
//	}
//	result = string(line[1:])
//	return
//
//}

func (c *Client) initConn(ctx context.Context, cn *pool.Conn) error {
	if cn.Inited {
		return nil
	}
	cn.Inited = true

	password := c.opt.Password
	//connPool := pool.NewSingleConnPool(c.connPool, cn)
	//cn
	//conn, err := c.connPool.Get(ctx)
	//if err != nil {
	//	return err
	//}

	cn.Write([]byte(fmt.Sprintf("auth %s\r\n", password)))
	cn.Write([]byte("ping\r\n"))
	reader := cn.Reader()
	defer cn.Clean()
	line, _, err := reader.ReadLine()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(string(line))
	return nil
}

func (c *Client) getConn(ctx context.Context) (*pool.Conn, error) {

	cn, err := c._getConn(ctx)
	if err != nil {
		return nil, err
	}

	return cn, nil
}

func (c *Client) _getConn(ctx context.Context) (*pool.Conn, error) {
	cn, err := c.connPool.Get(ctx)
	if err != nil {
		return nil, err
	}

	if cn.Inited {
		return cn, nil
	}

	if err := c.initConn(ctx, cn); err != nil {
		c.connPool.Remove(ctx, cn, err)
		if err := errors.Unwrap(err); err != nil {
			return nil, err
		}
		return nil, err
	}

	return cn, nil
}

func (c *Client) context(ctx context.Context) context.Context {
	if c.opt.ContextTimeoutEnabled {
		return ctx
	}
	return context.Background()
}
