package main

import (
	"fakebilibili/adapter/http/router"
	"fakebilibili/domain/servicce/contribution/socket"
)

func main() {
	go socket.Severe.Start()
	router.InitRouter()
}
