package router

import (
	"Search-Engine/web/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	engine := gin.Default()

	//engine.Use()
	engine.Use(middleware.CorsFunc)
	group := engine.Group("/api")
	{
		InitBaseRouter(group)
		InitIndexRouter(group)
		InitDataBaseRouter(group)
	}
	return engine
}
