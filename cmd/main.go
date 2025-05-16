package main

import (
	"fakebilibili/adapter/http/router"
	_ "fakebilibili/adapter/socket"
	"fakebilibili/crons"
	_ "fakebilibili/infrastructure/pkg/database"
	_ "fakebilibili/infrastructure/pkg/database/mysql"
	_ "fakebilibili/infrastructure/pkg/database/redis"
)

func main() {
	//这里还包括了消费者函数的启动
	crons.InitCron()
	router.InitRouter()
}
