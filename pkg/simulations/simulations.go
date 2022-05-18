package simulations

import (
	"github.com/kpetremann/claw-network/configs"
	"github.com/kpetremann/claw-network/pkg/topology"
)

func AssetsDown(graph *topology.Graph, devicesDown []string, linksDown []string) ([]string, error) {
	for _, device := range devicesDown {
		graph.Nodes[device].Status = false
	}

	for _, link := range linksDown {
		graph.Links[link].Status = false
	}

	defer graph.FullReset()
	graph.ComputeAllLinkStatus()

	isolated, err := graph.GetIsolatedBottomNodes()
	if err != nil {
		return nil, err
	}

	return isolated, nil
}

func RunAllNodesScenarios(graph *topology.Graph) (map[string][]string, error) {
	scenarios := make(map[string][]string)
	for _, node := range graph.Nodes {
		if node.Role == configs.BottomDeviceRole {
			continue
		}
		result, err := AssetsDown(graph, []string{node.Hostname}, nil)
		scenarios[node.Hostname] = result

		if err != nil {
			return nil, err
		}
	}

	return scenarios, nil
}
