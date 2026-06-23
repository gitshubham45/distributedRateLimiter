package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StrategyConfig struct {
	Name     string `json:"name"`
	Limit    int64  `json:"limit"`
	Interval int64  `json:"interval"`
}

type EndpointRegisterReq struct {
	Path     string         `json:"path"`
	Method   string         `json:"method"`
	Strategy StrategyConfig `json:"strategy"`
}

func RegisterEndpoint(c *gin.Context) {

	var req EndpointRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	switch req.Strategy.Name {

	case "per_hour":
	

		
	}

	// TODO: Add endpoint to database
	// TODO: Add strategy to cache

	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
