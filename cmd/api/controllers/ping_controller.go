package controllers

import (
	"github.com/gin-gonic/gin"
)

type PingController interface {
	Ping(context *gin.Context)
}

type pingControllerImpl struct{}

func NewPingController() PingController {
	return &pingControllerImpl{}
}

func (pingCtrl *pingControllerImpl) Ping(c *gin.Context) {

	c.JSON(200, gin.H{
		"code":         200,
		"message":      "pong",
	})
}
