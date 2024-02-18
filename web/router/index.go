package router

import (
	"Search-Engine/web/controller"
	"github.com/gin-gonic/gin"
)

func InitIndexRouter(group *gin.RouterGroup) {
	indexGroup := group.Group("")
	{
		indexGroup.POST("/addIndex", controller.AddIndex)
	}
}
