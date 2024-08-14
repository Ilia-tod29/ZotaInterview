package main

import (
	"ZotaInterview/api"
	"ZotaInterview/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

	server, err := api.NewServer(config)
	if err != nil {
		log.Fatal("cannot initiate the server", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
