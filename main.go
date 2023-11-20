package main

import (
	"jrsfun-server-go/cache"
	"jrsfun-server-go/manager"
	"jrsfun-server-go/router"
	"log"
)

func main() {
	manager.InitManager()
	cache.InitGoCaches()
	engine := router.StartRouter()
	err := engine.Run(":3000")
	if err != nil {
		log.Fatalf("run err:%v", err)
		return
	}
}
