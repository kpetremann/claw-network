package api

import (
	"fmt"

	"github.com/kpetremann/claw-network/configs"

	"github.com/gin-gonic/gin"
)

func ListenAndServer() {
	router := gin.Default()

	router.GET("/topology", ListTopology)

	router.POST("/topology/:topology", AddTopology)
	router.POST("/topology/custom/device/every/down/impact", SimulateDownImpactProvidedTopology)
	router.GET("/topology/:topology/device/every/down/impact", SimulateDownImpactExistingTopology)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	listenAddrPort := fmt.Sprintf("%s:%s", configs.ListenAddress, configs.ListenPort)

	router.Run(listenAddrPort)
}
