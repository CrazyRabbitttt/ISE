package router

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	engine := gin.Default()

	// 中间件处理跨域、异常
	//engine.Use()
	group1 := engine.Group("/api")
	{
		InitBaseRouter(group1)
		InitIndexRouter(group1)
		InitDataBaseRouter(group1)
	}
	return engine
}
