package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Unisponse(c *gin.Context, value int64) {
	type response struct {
		Value int32 `json:"value"`
	}
	c.JSON(http.StatusOK, gin.H{
		"value": value,
	})
}

func UnisponseBA(c *gin.Context, before, after int64) {
	type response struct {
		Previous int32 `json:"previous"`
		Value    int32 `json:"value"`
	}
	c.JSON(http.StatusOK, gin.H{
		"previous": before,
		"value":    after,
	})
}
