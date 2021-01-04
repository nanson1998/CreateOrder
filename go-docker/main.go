package main

import (
	"github.com/callicoder/go-docker/router"

	"log"
)

func main() {
	//redis.ConnectRd()
	//mysql.Connect()
	r := router.SetupRouter()
	log.Fatal(r.Run(":8081"))

}
