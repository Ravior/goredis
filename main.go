package main

import (
	"fmt"
	"github.com/Ravior/goredis/config"
	"github.com/Ravior/goredis/lib/logger"
	RedisServer "github.com/Ravior/goredis/redis/server"
	"github.com/Ravior/goredis/tcp"
	"os"
)

var banner = `
      ___           ___           ___                       ___     
     /\  \         /\  \         /\  \          ___        /\  \    
    /::\  \       /::\  \       /::\  \        /\  \      /::\  \   
   /:/\:\  \     /:/\:\  \     /:/\:\  \       \:\  \    /:/\ \  \  
  /::\~\:\  \   /::\~\:\  \   /:/  \:\__\      /::\__\  _\:\~\ \  \ 
 /:/\:\ \:\__\ /:/\:\ \:\__\ /:/__/ \:|__|  __/:/\/__/ /\ \:\ \ \__\
 \/_|::\/:/  / \:\~\:\ \/__/ \:\  \ /:/  / /\/:/  /    \:\ \:\ \/__/
    |:|::/  /   \:\ \:\__\    \:\  /:/  /  \::/__/      \:\ \:\__\  
    |:|\/__/     \:\ \/__/     \:\/:/  /    \:\__\       \:\/:/  /  
    |:|  |        \:\__\        \::/__/      \/__/        \::/  /   
     \|__|         \/__/         ~~                        \/__/    
`

var defaultProperties = &config.ServerProperties{
	Bind:           "0.0.0.0",
	Port:           6399,
	AppendOnly:     false,
	AppendFilename: "",
	MaxClients:     1000,
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}

func main() {
	print(banner)
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "go-redis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})
	configFileName := os.Getenv("CONFIG")
	if configFileName == "" {
		if fileExists("redis.conf") {
			config.SetupConfig("redis.conf")
		} else {
			config.Properties = defaultProperties
		}
	} else {
		config.SetupConfig(configFileName)
	}

	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, RedisServer.NewHandler())
	if err != nil {
		logger.Error(err)
	}
}
