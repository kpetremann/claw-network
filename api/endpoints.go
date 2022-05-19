package api

import (
	"github.com/kpetremann/claw-network/pkg/simulations"
	"github.com/kpetremann/claw-network/pkg/topology"

	"github.com/gin-gonic/gin"
)

func (s *SimulationManager) SimulateDownImpactExistingTopology(context *gin.Context) {
	var err error
	topologyName := context.Param("topology")
	deviceDown := context.Param("device")

	repo := <-s.getRepository

	graph, err := repo.LoadTopology(topologyName)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var result *simulations.SimulationResult
	if deviceDown == "each" {
		result, err = simulations.RunAllNodesScenarios(graph)
	} else {
		result, err = simulations.RunWithAssetsDown(graph, []string{deviceDown}, nil)
	}

	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.PureJSON(200, result)
}

func (s *SimulationManager) SimulateDownImpactProvidedTopology(context *gin.Context) {
	var scenarioParameters struct {
		graph       topology.Graph
		downDevices []string
		downLinks   []string
	}

	err := context.ShouldBind(&scenarioParameters)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	result, err := simulations.RunAllNodesScenarios(&scenarioParameters.graph)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, result)
}

func (s *SimulationManager) AddTopology(context *gin.Context) {
	topologyName := context.Param("topology")

	var graph topology.Graph
	if err := context.ShouldBind(&graph); err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	repo := <-s.getRepository
	if err := repo.SaveTopology(topologyName, &graph); err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	s.writeRepository <- repo.Topologies

	context.JSON(200, gin.H{"result": "topology saved"})
}

func (s *SimulationManager) ListTopology(context *gin.Context) {
	context.JSON(200, <-s.getRepository)
}

func (s *SimulationManager) GetTopology(context *gin.Context) {
	topologyName := context.Param("topology")
	repo := <-s.getRepository

	topo, err := repo.LoadTopology(topologyName)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, topo)
}

func (s *SimulationManager) DeleteTopology(context *gin.Context) {
	topologyName := context.Param("topology")
	repo := <-s.getRepository

	if err := repo.DeleteTopology(topologyName); err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, gin.H{"result": "deleted"})
}
