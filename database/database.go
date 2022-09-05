package database

import (
	"github.com/Ravior/goredis/interface/redis"
	"github.com/Ravior/goredis/redis/protocol"
	"sync/atomic"
)

// MultiDB is a set of multiple database set
type MultiDB struct {
	dbSet []*atomic.Value // *DB
}

// NewStandaloneServer creates a standalone redis server, with multi database and all other funtions
func NewStandaloneServer() *MultiDB {
	return &MultiDB{}
}

// Exec executes command
// parameter `cmdLine` contains command and its arguments, for example: "set key value"
func (mdb *MultiDB) Exec(conn redis.Connection, cmdLine [][]byte) (result redis.Reply) {
	// normal commands
	dbIndex := conn.GetDBIndex()
	selectedDB, errReply := mdb.selectDB(dbIndex)
	if errReply != nil {
		return errReply
	}
	return selectedDB.Exec(conn, cmdLine)
}

// AfterClientClose does some clean after client close connection
func (mdb *MultiDB) AfterClientClose(conn redis.Connection) {

}

func (mdb *MultiDB) Close() {

}

func (mdb *MultiDB) selectDB(dbIndex int) (*DB, *protocol.StandardErrReply) {
	if dbIndex >= len(mdb.dbSet) || dbIndex < 0 {
		return nil, protocol.NewErrReply("ERR DB index is out of range")
	}
	return mdb.dbSet[dbIndex].Load().(*DB), nil
}
