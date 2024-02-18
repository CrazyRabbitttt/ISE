package controller

import (
	"Search-Engine/search-engine/container"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, ResponseOkWithMessage("Welcome to use this search-engine"))
}

func Query(c *gin.Context) {
	c.JSON(http.StatusOK, ResponseOkWithMessage("Access the query function"))
}

func Cut(c *gin.Context) {
	query := c.Query("q")
	tokenizer := container.GlobalContainer.Tokenizer
	terms := tokenizer.Cut(query)
	c.JSON(http.StatusOK, ResponseOkWithMessage(terms))
}
