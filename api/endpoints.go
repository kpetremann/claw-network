package api

import (
	"fmt"
	"time"

	"github.com/kpetremann/claw-network/pkg/simulations"
	"github.com/kpetremann/claw-network/pkg/topology"

	"github.com/gin-gonic/gin"
)

func runAllScenarios(graph *topology.Graph) (map[string]interface{}, error) {
	start := time.Now()
	scenarios, err := simulations.RunAllNodesScenarios(graph)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	elapsed := t.Sub(start)

	msg := fmt.Sprintf("took %d ms", elapsed.Milliseconds())

	result := make(map[string]interface{})
	result["performance"] = msg
	result["impact_simulation"] = scenarios

	return result, err
}

func (s *SimulationManager) SimulateDownImpactExistingTopology(context *gin.Context) {
	topologyName := context.Param("topology")
	repo := <-s.getRepository

	graph, err := repo.LoadTopology(topologyName)

	if err != nil {
		context.JSON(500, err)
		return
	}

	result, err := runAllScenarios(graph)
	if err != nil {
		context.JSON(500, err)
		return
	}

	context.JSON(200, result)
}

func (s *SimulationManager) SimulateDownImpactProvidedTopology(context *gin.Context) {
	var graph topology.Graph
	err := context.ShouldBind(&graph)
	if err != nil {
		context.JSON(500, err)
		return
	}

	result, err := runAllScenarios(&graph)
	if err != nil {
		context.JSON(500, err)
		return
	}

	context.JSON(200, result)
}

func (s *SimulationManager) AddTopology(context *gin.Context) {
	topologyName := context.Param("topology")

	var graph topology.Graph
	err := context.ShouldBind(&graph)

	if err != nil {
		context.JSON(500, err)
		return
	}

	repo := <-s.getRepository
	if err := repo.SaveTopology(topologyName+".json", &graph); err != nil {
		context.JSON(500, err)
		return
	}

	context.JSON(200, "topology saved")
}

func (s *SimulationManager) ListTopology(context *gin.Context) {
	repo := <-s.getRepository

	if err := repo.UpdateTopology(); err != nil {
		context.JSON(500, err)
		return
	}

	context.JSON(200, repo)
}

func (s *SimulationManager) GetTopology(context *gin.Context) {
	context.JSON(200, "not implemented yet")
}

func (s *SimulationManager) DeleteTopology(context *gin.Context) {
	topologyName := context.Param("topology")
	repo := <-s.getRepository

	if err := repo.DeleteTopology(topologyName); err != nil {
		context.JSON(500, err)
		return
	}

	context.JSON(200, "deleted")
}
