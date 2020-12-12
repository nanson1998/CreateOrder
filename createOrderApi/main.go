package main

import (
	"Api/helper/mysql"
	"Api/helper/redis"
	"Api/router"

	"log"
)

func main() {
	redis.ConnectRd()
	mysql.Connect()
	r := router.SetupRouter()
	log.Fatal(r.Run(":8080"))

}
