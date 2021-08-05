package main

import (
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/mustafasegf/hompimpa/api"
	"github.com/mustafasegf/hompimpa/util"
)

func main() {
	err := util.SetLogger()
	if err != nil {
		log.Fatal("cannot set logger: ", err)
	}

	config, err := util.LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	pub := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		DB:   0,
	})

	sub := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		DB:   0,
	})

	server := api.MakeServer(config, pub, sub)
	server.RunServer()
}
