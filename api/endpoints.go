package api

import (
	"fmt"
	"time"

	"github.com/kpetremann/claw-network/configs"
	"github.com/kpetremann/claw-network/internal/backends"
	"github.com/kpetremann/claw-network/pkg/simulations"
	"github.com/kpetremann/claw-network/pkg/topology"

	"github.com/gin-gonic/gin"
)

func SimulationDownImpactExample(context *gin.Context) {
	graph, err := backends.LoadTopologyFromFile(configs.TopologyFile)

	if err != nil {
		context.JSON(500, err)
		return
	}

	start := time.Now()
	scenario := make(map[string][]string)
	for _, node := range graph.Nodes {
		if node.Role == configs.BottomDeviceRole {
			continue
		}

		if scenario[node.Hostname], err = simulations.AssetsDown(graph, []string{node.Hostname}, nil); err != nil {
			context.JSON(500, err)
			return
		}
	}
	t := time.Now()
	elapsed := t.Sub(start)

	msg := fmt.Sprintf("took %d ms", elapsed.Milliseconds())

	result := make(map[string]interface{})
	result["performances"] = msg
	result["impact_simulation"] = scenario

	context.JSON(200, result)
}

func SimulateDownImpactProvidedTopology(context *gin.Context) {
	var graph topology.Graph
	err := context.ShouldBind(&graph)

	if err != nil {
		context.JSON(500, err)
		return
	}

	scenarios := make(map[string][]string)
	for _, node := range graph.Nodes {
		if node.Role == configs.BottomDeviceRole {
			continue
		}
		if scenarios[node.Hostname], err = simulations.AssetsDown(&graph, []string{node.Hostname}, nil); err != nil {
			context.JSON(500, err)
			return
		}
	}

	context.JSON(200, scenarios)
}

func SimulateDownImpactExistingTopology(context *gin.Context) {
	topo := context.Param("topology")
	context.JSON(501, &topo)
}

func AddTopology(context *gin.Context) {
	topologyName := context.Param("topology")
	var graph topology.Graph
	err := context.ShouldBind(&graph)

	if err != nil {
		context.JSON(500, err)
		return
	}

	if err := backends.SaveTopologyToFile(topologyName+".json", &graph); err != nil {
		context.JSON(500, err)
		return
	}

	context.JSON(200, "topology saved")
}

func DeleteTopology(context *gin.Context) {
	context.JSON(501, "not implemented yet")
}

func ListTopology(context *gin.Context) {
	context.JSON(501, "not implemented yet")
}

func GetTopology(context *gin.Context) {
	context.JSON(200, "not implemented yet")
}
