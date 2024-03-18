package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CorsFunc(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Set("Access-Control-Max-Age", "86400")
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.Next()
}

func LogHttpRequest(c *gin.Context) {
	// 打印 HTTP 请求方法和路径
	fmt.Printf("收到请求：%s %s\n", c.Request.Method, c.Request.URL.Path)

	// 打印 HTTP 请求头
	fmt.Println("请求头：")
	for key, values := range c.Request.Header {
		fmt.Printf("%s: %s\n", key, values)
	}

	// 打印 HTTP 请求体
	body, _ := c.GetRawData()
	fmt.Println("请求体：", string(body), "-------")

	// 继续处理请求
	c.Next()
}
