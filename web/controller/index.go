package controller

import (
	"Search-Engine/search-engine/model"
	"Search-Engine/web/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddIndex(c *gin.Context) {
	indexDoc := &model.IndexDoc{}
	err := c.ShouldBind(&indexDoc)
	if err != nil {
		ResponseErrWithMessage("解析Http请求到结构体(indexDoc)失败")
	}
	fmt.Println("terms:", indexDoc.Text)
	fmt.Println("key:", indexDoc.Key)
	for k, v := range indexDoc.Attrs {
		fmt.Println("attr_k:", k, " attr_v:", v)
	}
	err = service.GlobalService.IndexService.AddIndexDoc(indexDoc)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseErrWithMessage("添加索引失败"))
	}
	c.JSON(http.StatusOK, ResponseOkWithMessage("添加索引成功"))
}
