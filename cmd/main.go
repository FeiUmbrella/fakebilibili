package main

import (
	"fakebilibili/adapter/http/router"
	_ "fakebilibili/adapter/socket"
	_ "fakebilibili/infrastructure/pkg/database"
	_ "fakebilibili/infrastructure/pkg/database/mysql"
	_ "fakebilibili/infrastructure/pkg/database/redis"
)

func main() {
	router.InitRouter()
}
