package contribution

import (
	"fakebilibili/adapter/http/controller"
	"github.com/gin-gonic/gin"
)

type Controllers struct {
	controller.BaseControllers
}

func (c *Controllers) GetVideoBarrage(ctx *gin.Context) {
	
}
