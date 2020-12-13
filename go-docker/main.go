package main

import (
	"github.com/callicoder/go-docker/helper/mysql"
	"github.com/callicoder/go-docker/helper/redis"
	"github.com/callicoder/go-docker/router"

	"log"
)

func main() {
	redis.ConnectRd()
	mysql.Connect()
	r := router.SetupRouter()
	log.Fatal(r.Run(":8081"))

}
