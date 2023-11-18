package main

import (
	"TogetherForever/config"
	"TogetherForever/gui"
	"TogetherForever/server"
)

var (
	Config = config.LoadConfig("config.toml")
)

func main() {
	if Config.Enable {
		gui.InitGui()
	} else {
		gui.Enabled = false
	}

	if gui.Enabled {
		server.InitServer(Config)
		go server.MainServer.Start()
		gui.StartGui()
	} else {
		server.InitServer(Config)
		server.MainServer.Start()
	}
}
