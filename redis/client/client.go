package client

import "net"

const (
	created = iota
	running
	closed
)

type Client struct {
	conn   net.Conn
	status int32
}
