package main

import (
	"net"
	"net/textproto"
)

type Conn struct {
	conn net.Conn
	text *textproto.Conn
}

func NewConn(c net.Conn) *Conn {
	return &Conn{
		conn: c,
		text: textproto.NewConn(c),
	}
}

// Close closes both the conn and the text connections
func (c *Conn) Close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}
	if err := c.text.Close(); err != nil {
		return err
	}
	return nil
}

// readLine reads a line from the connection
func (c *Conn) readLine() (string, error) {
	return c.text.ReadLine()
}

// writeResponse writes a line to the connection
func (c *Conn) writeResponse(code int, text string) error {
	return c.text.PrintfLine("%d %v", code, text)
}
