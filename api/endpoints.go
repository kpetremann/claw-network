package api

import (
	"fmt"
	"time"

	"github.com/kpetremann/claw-network/internal/backends"
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

func SimulateDownImpactProvidedTopology(context *gin.Context) {
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

func SimulateDownImpactExistingTopology(context *gin.Context) {
	var repo backends.TopologyRepository
	topologyName := context.Param("topology")
	graph, err := repo.LoadTopology(topologyName + ".json")

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

func AddTopology(context *gin.Context) {
	topologyName := context.Param("topology")
	var repo backends.TopologyRepository
	var graph topology.Graph
	err := context.ShouldBind(&graph)

	if err != nil {
		context.JSON(500, err)
		return
	}

	if err := repo.SaveTopology(topologyName+".json", &graph); err != nil {
		context.JSON(500, err)
		return
	}

	context.JSON(200, "topology saved")
}

func ListTopology(context *gin.Context) {
	var repo backends.TopologyRepository
	if err := repo.UpdateTopology(); err != nil {
		context.JSON(500, err)
		return
	}

	context.JSON(200, repo)
}

func GetTopology(context *gin.Context) {
	context.JSON(200, "not implemented yet")
}

func DeleteTopology(context *gin.Context) {
	topologyName := context.Param("topology")
	var repo backends.TopologyRepository

	if err := repo.DeleteTopology(topologyName); err != nil {
		context.JSON(500, err)
		return
	}

	context.JSON(200, "deleted")
}
