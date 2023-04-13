package pool

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"sync/atomic"
	"time"
)

var noDeadline = time.Time{}

type Conn struct {
	usedAt  int64 // atomic
	netConn net.Conn

	rd *bufio.Reader
	wr *bufio.Writer

	Inited    bool
	pooled    bool
	createdAt time.Time
}

func (c *Conn) Reader() *bufio.Reader {
	return c.rd
}
func (c *Conn) Writer() *bufio.Writer {
	return c.wr
}
func NewConn(netConn net.Conn) *Conn {
	cn := &Conn{
		netConn:   netConn,
		createdAt: time.Now(),
	}
	cn.rd = bufio.NewReader(netConn)
	cn.wr = bufio.NewWriter(netConn)
	cn.SetUsedAt(time.Now())
	return cn
}

func (cn *Conn) SetUsedAt(tm time.Time) {
	atomic.StoreInt64(&cn.usedAt, tm.Unix())
}

func (cn *Conn) UsedAt() time.Time {
	unix := atomic.LoadInt64(&cn.usedAt)
	return time.Unix(unix, 0)
}
func (cn *Conn) SetNetConn(netConn net.Conn) {
	cn.netConn = netConn
	cn.rd.Reset(netConn)
	cn.wr.Reset(netConn)
}
func (cn *Conn) Write(b []byte) (int, error) {
	return cn.netConn.Write(b)
}

func (cn *Conn) RemoteAddr() net.Addr {
	if cn.netConn != nil {
		return cn.netConn.RemoteAddr()
	}
	return nil
}
func (cn *Conn) WithReader(
	ctx context.Context, timeout time.Duration, fn func(rd *bufio.Reader) error,
) error {
	if timeout >= 0 {
		if err := cn.netConn.SetReadDeadline(cn.deadline(ctx, timeout)); err != nil {
			return err
		}
	}
	return fn(cn.rd)
}

func (cn *Conn) WithWriter(
	ctx context.Context, timeout time.Duration, fn func(wr *bufio.Writer) error,
) error {
	if timeout >= 0 {
		if err := cn.netConn.SetWriteDeadline(cn.deadline(ctx, timeout)); err != nil {
			return err
		}
	}

	if cn.wr.Buffered() > 0 {
		cn.wr.Reset(cn.netConn)
	}

	if err := fn(cn.wr); err != nil {
		return err
	}

	return cn.wr.Flush()
}

func (cn *Conn) Close() error {
	return cn.netConn.Close()
}

func (cn *Conn) deadline(ctx context.Context, timeout time.Duration) time.Time {
	tm := time.Now()
	cn.SetUsedAt(tm)

	if timeout > 0 {
		tm = tm.Add(timeout)
	}

	if ctx != nil {
		deadline, ok := ctx.Deadline()
		if ok {
			if timeout == 0 {
				return deadline
			}
			if deadline.Before(tm) {
				return deadline
			}
			return tm
		}
	}

	if timeout > 0 {
		return tm
	}

	return noDeadline
}

func (cn *Conn) Clean() {
	cn.rd.Reset(cn.netConn)
}
func (cn *Conn) Do(ctx context.Context, command string, options ...string) (result string, err error) {
	opt := strings.Builder{}
	for _, option := range options {
		opt.WriteByte(' ')
		opt.WriteString(option)
	}
	if err != nil {
		return "", err
	}
	rd, _ := cn.Reader(), cn.Writer()
	defer cn.Clean()

	_, err = cn.Write([]byte(fmt.Sprint(command, opt.String(), " \r\n")))
	if err != nil {
		return "", err
	}
	line, _, err := rd.ReadLine()
	fmt.Println(string(line))
	if err != nil {
		log.Fatalln(err)
	}
	switch line[0] {
	case '-':
		return "", errors.New(string(line[1:]))
	default:
		return string(line[1:]), nil
	}
}
