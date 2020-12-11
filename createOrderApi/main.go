package main

import (
	"Api/router"
	"log"
)

func main() {
	//redis.ConnectRd()
	r := router.SetupRouter()
	log.Fatal(r.Run(":8082"))

}
