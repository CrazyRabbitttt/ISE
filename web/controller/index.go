package controller

import (
	"Search-Engine/search-engine/model"
	"Search-Engine/web/service"
	"github.com/gin-gonic/gin"
)

func AddIndex(c *gin.Context) {
	indexDoc := &model.IndexDoc{}
	err := c.ShouldBind(&indexDoc)
	if err != nil {
		ResponseErrWithMessage("结构化读取Request请求失败, parse indexDoc")
	}
	err = service.GlobalService.IndexService.AddIndexDoc(indexDoc)
	if err != nil {
		ResponseErrWithMessage("添加索引失败")
	}
	ResponseOkWithMessage("添加索引成功")
}
