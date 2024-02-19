package router

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	engine := gin.Default()
	// 中间件处理跨域、异常

	//engine.Use()
	group := engine.Group("/api")
	{
		InitBaseRouter(group)
		InitIndexRouter(group)
		InitDataBaseRouter(group)
	}
	return engine
}
