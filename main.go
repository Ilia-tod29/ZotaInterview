package main

import (
	"ZotaInterview/api"
	"ZotaInterview/client"
	"ZotaInterview/util"
	"log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}
	c := client.GetZotaClient(config.SecretKey)

	server, err := api.NewServer(config, c)
	if err != nil {
		log.Fatal("cannot initiate the server", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
