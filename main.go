package main

import (
	"golife/config"
	"golife/web"
	"log"
)

func main() {
	// load toml config
	var conf config.Config
	if err := conf.LoadFile("./config.toml"); err != nil {
		log.Fatal(err)
	}
	// create & start server
	server := web.Server{Config: conf}
	log.Fatal(server.Start())
}
