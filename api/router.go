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

	router.GET("/topology/:topology", s.GetTopology)
	router.POST("/topology/:topology", s.AddTopology)
	router.DELETE("/topology/:topology", s.DeleteTopology)

	router.GET("/topology/:topology/anomalies", func(c *gin.Context) { c.JSON(501, gin.H{"error": "not implemented"}) })

	router.POST("/topology/custom/device/:device/down/impact", s.RunOnProvidedTopology)
	router.GET("/topology/:topology/device/:device/down/impact", s.RunOnExistingTopology)

	router.POST("/topology/custom/link/:link/down/impact", s.RunOnProvidedTopology)
	router.GET("/topology/:topology/link/:link/down/impact", s.RunOnExistingTopology)

	router.POST("/topology/:topology/scenario/impact", s.RunScenario)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	listenAddrPort := fmt.Sprintf("%s:%s", configs.ListenAddress, configs.ListenPort)

	router.Run(listenAddrPort)
}
