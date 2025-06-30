package routers

import (
	"github.com/gin-gonic/gin"
	"zg6zy5/apiway/handler/apihandler"
)

func LoadRouters() {
	router := gin.Default()

	// 简单的路由组: v1
	{
		v1 := router.Group("/v1")

		v1.GET("/sign", apihandler.Sign)
		v1.GET("/one", apihandler.One)
		v1.GET("/calblack", apihandler.Calblack)
	}

	router.Run(":8080")
}
