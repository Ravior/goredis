package database

import "strings"

var cmdTable = make(map[string]*command)

type command struct {
	executor ExecFunc
	prepare  PreFunc // return related keys command
}

// RegisterCommand registers a new command
func RegisterCommand(name string, executor ExecFunc) {
	name = strings.ToLower(name)
	cmdTable[name] = &command{
		executor: executor,
	}
}
