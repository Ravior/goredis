package database

import (
	"github.com/Ravior/goredis/interface/redis"
	"github.com/Ravior/goredis/redis/protocol"
)

// execDel removes a key from db
func execDel(db *DB, args [][]byte) redis.Reply {
	keys := make([]string, len(args))
	for i, v := range args {
		keys[i] = string(v)
	}

	deleted := db.Removes(keys...)
	if deleted > 0 {
		//db.addAof(utils.ToCmdLine3("del", args...))
	}
	return protocol.NewIntReply(int64(deleted))
}

func init() {
	RegisterCommand("Del", execDel)
}
