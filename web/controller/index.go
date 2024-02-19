package controller

import (
	"Search-Engine/search-engine/model"
	"Search-Engine/web/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AddIndex(c *gin.Context) {
	indexDoc := &model.IndexDoc{}
	err := c.ShouldBind(&indexDoc)
	if err != nil {
		ResponseErrWithMessage("解析Http请求到结构体(indexDoc)失败")
	}
	err = service.GlobalService.IndexService.AddIndexDoc(indexDoc)
	//fmt.Println("IndexService处理完IndexDoc")
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseErrWithMessage("添加索引失败"))
	}
	c.JSON(http.StatusOK, ResponseOkWithMessage("添加索引成功"))
}
