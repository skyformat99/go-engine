package conn

import (
	"context"
	"errors"
	"net"
)

type tcpConn struct {
	conn     *net.TCPConn
	listener *net.TCPListener
	cancel   context.CancelFunc
	info     string
}

func (c *tcpConn) Name() string {
	return "tcp"
}

func (c *tcpConn) Read(p []byte) (n int, err error) {
	if c.conn != nil {
		return c.conn.Read(p)
	}
	return 0, errors.New("empty conn")
}

func (c *tcpConn) Write(p []byte) (n int, err error) {
	if c.conn != nil {
		return c.conn.Write(p)
	}
	return 0, errors.New("empty conn")
}

func (c *tcpConn) Close() error {
	if c.cancel != nil {
		c.cancel()
	}
	if c.conn != nil {
		return c.conn.Close()
	} else if c.listener != nil {
		return c.listener.Close()
	}
	return nil
}

func (c *tcpConn) Info() string {
	if c.info != "" {
		return c.info
	}
	if c.conn != nil {
		c.info = c.conn.LocalAddr().String() + "<--tcp-->" + c.conn.RemoteAddr().String()
	} else if c.listener != nil {
		c.info = "tcp--" + c.listener.Addr().String()
	} else {
		c.info = "empty tcp conn"
	}
	return c.info
}

func (c *tcpConn) Dial(dst string) (Conn, error) {
	addr, err := net.ResolveTCPAddr("tcp", dst)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", addr.String())
	if err != nil {
		return nil, err
	}
	c.cancel = nil
	return &tcpConn{conn: conn.(*net.TCPConn)}, nil
}

func (c *tcpConn) Listen(dst string) (Conn, error) {
	addr, err := net.ResolveTCPAddr("tcp", dst)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &tcpConn{listener: listener}, nil
}

func (c *tcpConn) Accept() (Conn, error) {
	conn, err := c.listener.Accept()
	if err != nil {
		return nil, err
	}
	return &tcpConn{conn: conn.(*net.TCPConn)}, nil
}
