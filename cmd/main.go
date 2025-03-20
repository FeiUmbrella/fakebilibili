package main

import (
	"fakebilibili/infrastructure/pkg/global"
	"fmt"
)

func main() {
	sqlDb, _ := global.Db.DB()
	err := sqlDb.Ping()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Success")
}
