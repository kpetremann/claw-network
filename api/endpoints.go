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

	graph, err := s.repository.LoadTopology(topologyName)

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

func (s *SimulationManager) RunScenario(context *gin.Context) {
	var params struct {
		Graph       interface{} `json:"topology"`
		DevicesDown []string    `json:"devices_down"`
		LinksDown   []string    `json:"links_down"`
	}

	err := context.ShouldBind(&params)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	var graph *topology.Graph
	switch params.Graph.(type) {
	case *topology.Graph:
		graph = params.Graph.(*topology.Graph)
	case string:
		graph, err = s.repository.LoadTopology(params.Graph.(string))
		if err != nil {
			context.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	result, err := simulations.RunWithAssetsDown(graph, params.DevicesDown, params.LinksDown)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, result)
}

func (s *SimulationManager) GetAnomalies(context *gin.Context) {
	topologyName := context.Param("topology")

	topo, err := s.repository.LoadTopology(topologyName)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	anomalies := topo.GetAnomalies()
	context.JSON(200, anomalies)
}

func (s *SimulationManager) AddTopology(context *gin.Context) {
	topologyName := context.Param("topology")

	var graph topology.Graph
	if err := context.ShouldBind(&graph); err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if err := s.repository.SaveTopology(topologyName, &graph); err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, gin.H{"result": "saved"})
}

func (s *SimulationManager) ListTopology(context *gin.Context) {
	context.JSON(200, s.repository.GetTopologies())
}

func (s *SimulationManager) ListTopologiesDetails(context *gin.Context) {
	topoListDetails, err := s.repository.ListTopologiesDetails()

	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, topoListDetails)
}

func (s *SimulationManager) GetTopology(context *gin.Context) {
	topologyName := context.Param("topology")

	topo, err := s.repository.LoadTopology(topologyName)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, topo)
}

func (s *SimulationManager) GetTopologyDetails(context *gin.Context) {
	topologyName := context.Param("topology")

	topo, err := s.repository.GetTopologyDetails(topologyName)
	if err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, topo)
}

func (s *SimulationManager) DeleteTopology(context *gin.Context) {
	topologyName := context.Param("topology")

	if err := s.repository.DeleteTopology(topologyName); err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}

	context.JSON(200, gin.H{"result": "deleted"})
}
