package controller

import (
	"Search-Engine/search-engine/model"
	"Search-Engine/web/service"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"net/http"
)

func DebugIndex(c *gin.Context) {
	indexDoc := &model.IndexDoc{}
	err := c.ShouldBind(&indexDoc)
	if err != nil {
		ResponseErrWithMessage("解析Http请求到结构体(indexDoc)失败")
	}
	// 打印将收到的HTTP请求绑定到 indexDoc 后的各个属性值
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println("Error initializing snowflake node:", err)
		return
	}
	uuid := node.Generate()
	indexDoc.Key = uuid.Int64()
	fmt.Printf("addIndex, docId:%d, terms:%s, title:%s, url:%s\n", indexDoc.Key,
		indexDoc.Text, indexDoc.Attrs["title"], indexDoc.Attrs["page_url"])
	service.GlobalService.IndexService.DebugPositiveIndex(indexDoc)
	c.JSON(http.StatusOK, ResponseOkWithMessage("debuging......."))
}

func AddIndex(c *gin.Context) {
	indexDoc := &model.IndexDoc{}
	err := c.ShouldBind(&indexDoc)
	if err != nil {
		ResponseErrWithMessage("解析Http请求到结构体(indexDoc)失败")
	}
	// 打印将收到的HTTP请求绑定到 indexDoc 后的各个属性值
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println("Error initializing snowflake node:", err)
		return
	}
	uuid := node.Generate()
	indexDoc.Key = uuid.Int64()
	fmt.Printf("addIndex, docId:%d, terms:%s, title:%s, url:%s, description:%s|||\n", indexDoc.Key,
		indexDoc.Text, indexDoc.Attrs["title"], indexDoc.Attrs["page_url"], indexDoc.Attrs["description"])
	err = service.GlobalService.IndexService.AddIndexDoc(indexDoc)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseErrWithMessage("添加索引失败"))
	}
	c.JSON(http.StatusOK, ResponseOkWithMessage("添加索引成功"))
}
