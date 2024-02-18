package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, ResponseOkWithMessage("Welcome to use this search-engine"))
}

func Query(c *gin.Context) {
	c.JSON(http.StatusOK, ResponseOkWithMessage("Access the query function"))
}
