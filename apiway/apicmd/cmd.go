package main

import (
	"github.com/gin-gonic/gin"
	"zg6zy5/apiway/inits"
	"zg6zy5/apiway/routers"
)

func main() {
	inits.ExampleClient()
	inits.InitMysql()
	router := gin.Default()

	routers.LoadRouters()

	router.GET("/v1/sign", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // 监听并在 0.0.0.0:8080 上启动服务
}
