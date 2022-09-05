package database

import "github.com/Ravior/goredis/interface/redis"

type DB interface {
	Exec(client redis.Connection, cmdLine [][]byte) redis.Reply
	AfterClientClose(conn redis.Connection)
	Close()
}
