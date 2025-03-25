package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Cors 跨域中间件，如果前后端部署在不同IP地址，前端访问后端要进行跨域配置
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求方法
		method := c.Request.Method
		// 请求头部
		origin := c.Request.Header.Get("Origin")
		// 声明请求头keys
		var headerKeys []string
		for k := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			//这是允许访问所有域
			c.Header("Access-Control-Allow-Origin", "*")
			//服务器支持的所有跨域请求的方法,为了避免浏览次请求的多次'预检'请求
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//header的类型
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session,X_Requested_With,Accept, Origin, Host, Connection, Accept-Encoding, Accept-Language,DNT, X-CustomHeader, Keep-Alive, User-Agent, X-Requested-With, If-Modified-Since, Cache-Control, Content-Type, Pragma")
			//允许跨域设置可以返回其他子段  跨域关键设置 让浏览器可以解析
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
			//缓存预检请求信息的时间(OPTIONS) 单位为秒 2天
			c.Header("Access-Control-Max-Age", "172800")
			//跨域请求是否需要带cookie信息 默认设置为true 如果允许携带cookie 设置为true Allow-Origin 不能为*
			c.Header("Access-Control-Allow-Credentials", "false")
			//设置返回格式是json
			c.Set("content-type", "application/json")
		}

		//放行所有OPTIONS方法, 遇到OPTIONS请求直接返回响应，防止进入业务逻辑
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		// 让出控制权，处理后序中间件或路由
		c.Next()
	}
}
