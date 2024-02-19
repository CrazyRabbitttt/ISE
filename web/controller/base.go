package controller

import (
	"Search-Engine/search-engine/container"
	"Search-Engine/search-engine/model"
	"Search-Engine/web/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, ResponseOkWithMessage("Welcome to use this search-engine"))
}

func Query(c *gin.Context) {
	queryRequest := &model.SearchRequest{}
	if err := c.ShouldBind(&queryRequest); err != nil {
		c.JSON(http.StatusBadRequest, ResponseErrWithMessage("解析Http请求到结构体(SearchRequest)失败"))
	}
	service.GlobalService.BaseService.Query(queryRequest)
	c.JSON(http.StatusOK, ResponseOkWithMessage("Access the query function"))
}

func Cut(c *gin.Context) {
	query := c.Query("q")
	tokenizer := container.GlobalContainer.Tokenizer
	terms := tokenizer.Cut(query)
	c.JSON(http.StatusOK, ResponseOkWithMessage(terms))
}
