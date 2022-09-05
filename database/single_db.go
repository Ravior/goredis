package database

import (
	"github.com/Ravior/goredis/datastruct/dict"
	"github.com/Ravior/goredis/interface/redis"
	"github.com/Ravior/goredis/redis/protocol"
	"strings"
)

type DB struct {
	index int
	// key -> DataEntity
	data dict.Dict
}

// ExecFunc is interface for command executor
// args don't include cmd line
type ExecFunc func(db *DB, args [][]byte) redis.Reply

// PreFunc analyses command line when queued command to `multi`
// returns related write keys and read keys
type PreFunc func(args [][]byte) ([]string, []string)

// Exec executes command within one database
func (db *DB) Exec(c redis.Connection, cmdLine [][]byte) redis.Reply {
	// transaction control commands and other commands which cannot execute within transaction
	//cmdName := strings.ToLower(string(cmdLine[0]))
	return db.execNormalCommand(cmdLine)
}

func (db *DB) execNormalCommand(cmdLine [][]byte) redis.Reply {
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return protocol.NewErrReply("ERR unknown command '" + cmdName + "'")
	}
	fun := cmd.executor
	return fun(db, cmdLine[1:])
}

// Remove the given key from db
func (db *DB) Remove(key string) {
}

// Removes the given keys from db
func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}
