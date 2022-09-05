package connection

import (
	"github.com/Ravior/goredis/lib/sync/wait"
	"net"
	"sync"
	"time"
)

type Connection struct {
	conn net.Conn

	// waiting until protocol finished
	waitingReply wait.Wait

	// lock while server sending response
	mu sync.RWMutex

	// password may be changed by CONFIG command during runtime, so store the password
	password string

	// selected db
	selectedDB int
}

// NewConn creates Connection instance
func NewConn(conn net.Conn) *Connection {
	return &Connection{
		conn: conn,
	}
}

// RemoteAddr returns the remote network address
func (c *Connection) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

func (c *Connection) SetPassword(password string) {
	c.password = password
}

func (c *Connection) GetPassword() string {
	return c.password
}

// GetDBIndex returns selected db
func (c *Connection) GetDBIndex() int {
	return c.selectedDB
}

// SelectDB selects a database
func (c *Connection) SelectDB(dbIndex int) {
	c.selectedDB = dbIndex
}

// Write sends response to client over tcp connection
func (c *Connection) Write(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	c.waitingReply.Add(1)
	defer func() {
		c.waitingReply.Done()
	}()

	_, err := c.conn.Write(b)
	return err
}

// Close disconnect with the client
func (c *Connection) Close() error {
	c.waitingReply.WaitWithTimeout(10 * time.Second)
	_ = c.conn.Close()
	return nil
}
