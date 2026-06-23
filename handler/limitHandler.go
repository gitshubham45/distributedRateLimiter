package handler

import (
	"net/http"
	"rate-limiter/limiter"

	"github.com/gin-gonic/gin"
)

var strategyMap map[string]string

func ChooseStrategy(path string, strategy string) {

}



func GetLimiter(path string) RequestLimiter {
	return &limiter.PerHourLimit{}
}

func HandleLimit(c *gin.Context) {

	path := c.Request.URL.Path

	requestLimiter := GetLimiter(path)

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
