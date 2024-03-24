package controller

import (
	"Search-Engine/search-engine/container"
	"Search-Engine/search-engine/model"
	"Search-Engine/web/service"
	"fmt"
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
		return
	}
	//fmt.Println("解析Request请求结构成功, q:", queryRequest.Query, "limit:", queryRequest.Limit)
	response, err := service.GlobalService.BaseService.Query(queryRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseErrWithMessage("查询出现异常"))
	}
	c.JSON(http.StatusOK, ResponseOkWithMessage(response))
}

func Cut(c *gin.Context) {
	query := c.Query("q")
	tokenizer := container.GlobalContainer.Tokenizer
	terms := tokenizer.Cut(query)
	c.JSON(http.StatusOK, ResponseOkWithMessage(terms))
}

func SearchRemind(c *gin.Context) {
	query := c.Query("q")
	fmt.Printf("for remind, the search query is %s", query)
	res, err := service.GlobalService.BaseService.SearchRemind(query)
	fmt.Println("The len of remind words:", len(res))
	if err != nil {
		c.JSON(http.StatusBadRequest, ResponseErrWithMessage("搜索信息提示处理异常"))
	} else {
		c.JSON(http.StatusOK, ResponseOkWithMessage(res))
	}
}

func InitReminder(c *gin.Context) {
	request := &model.IndexDoc{}
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, ResponseErrWithMessage("parse model error"))
	}
	var querys []string
	for _, v := range request.Attrs {
		querys = append(querys, v)
	}
	fmt.Println("len of request querys:", len(querys))
	fmt.Printf("before init, print querys:%v", querys)
	service.GlobalService.BaseService.InitReminder(querys)
	c.JSON(http.StatusOK, ResponseOkWithMessage("init reminder success"))
}
