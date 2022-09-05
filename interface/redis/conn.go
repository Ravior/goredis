package redis

// Connection represents a connection with redis client
type Connection interface {
	Write([]byte) error
	SetPassword(string)
	GetPassword() string

	GetDBIndex() int
	SelectDB(int)
}
