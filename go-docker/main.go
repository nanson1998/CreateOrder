package main

import (
	"log"
	"zalopay-api/router"
)

func main() {
	r := router.SetupRouter()
	log.Fatal(r.Run(":9831"))
}
