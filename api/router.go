package api

import (
	"fmt"

	"github.com/kpetremann/claw-network/configs"

	"github.com/gin-gonic/gin"
)

// https://github.com/eliben/code-for-blog/blob/master/2019/gohttpconcurrency/channel-manager-server.go

func ListenAndServer() {
	router := gin.Default()
	s := NewSimulationManager()

	router.GET("/topology", s.ListTopology)

	router.POST("/topology/:topology", s.AddTopology)
	router.DELETE("/topology/:topology", s.DeleteTopology)
	router.POST("/topology/custom/device/each/down/impact", s.SimulateDownImpactProvidedTopology)
	router.GET("/topology/:topology/device/each/down/impact", s.SimulateDownImpactExistingTopology)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	listenAddrPort := fmt.Sprintf("%s:%s", configs.ListenAddress, configs.ListenPort)

	router.Run(listenAddrPort)
}
