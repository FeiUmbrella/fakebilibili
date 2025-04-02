package socket

import "fakebilibili/domain/servicce/users"

func init() {
	//初始化所有socket
	go users.Severe.Start()
}
