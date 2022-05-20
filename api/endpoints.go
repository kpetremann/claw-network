package api

import (
	"github.com/kpetremann/claw-network/pkg/simulations"
	"github.com/kpetremann/claw-network/pkg/topology"

	"github.com/gin-gonic/gin"
)

func (s *SimulationManager) RunOnExistingTopology(context *gin.Context) {
	var err error

	topologyName := context.Param("topology")
	deviceDown := context.Param("device")
	linkDown := context.Param("link")

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
		result, err = simulations.RunWithAssetsDown(graph, []string{deviceDown}, []string{linkDown})
	}

	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, result)
}

func (s *SimulationManager) RunOnProvidedTopology(context *gin.Context) {
	var params struct {
		Graph *topology.Graph `json:"topology"`
	}

	err := context.ShouldBind(&params)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	deviceDown := context.Param("device")
	linkDown := context.Param("link")

	var result *simulations.SimulationResult
	switch {
	case deviceDown == "each":
		result, err = simulations.RunAllNodesScenarios(params.Graph)
	default:
		result, err = simulations.RunWithAssetsDown(params.Graph, []string{deviceDown}, []string{linkDown})
	}
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
