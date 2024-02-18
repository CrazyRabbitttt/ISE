package router

import (
	"Search-Engine/web/controller"
	"github.com/gin-gonic/gin"
)

func InitBaseRouter(group *gin.RouterGroup) {
	baseRouter := group.Group("")
	{
		baseRouter.GET("/ping", controller.Welcome)
		baseRouter.GET("/cut", controller.Cut)
		baseRouter.POST("/query", controller.Query)
	}
}
