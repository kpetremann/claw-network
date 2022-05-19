package simulations

import (
	"fmt"

	"github.com/kpetremann/claw-network/configs"
	"github.com/kpetremann/claw-network/internal/utils"
	"github.com/kpetremann/claw-network/pkg/topology"
)

type ScenarioParameters struct {
	DevicesDown []string `json:"devices_down"`
	LinksDown   []string `json:"links_down"`
}

type ScenarioResult struct {
	ImpactedNodes []string           `json:"impacts"`
	Parameters    ScenarioParameters `json:"parameters"`
}

type SimulationResult struct {
	ScenarioResults map[string]*ScenarioResult `json:"scenarios_result"`
	Performance     string                     `json:"compute_time"`
}

func runScenario(graph *topology.Graph, devicesDown []string, linksDown []string) (*ScenarioResult, error) {
	result := ScenarioResult{
		Parameters: ScenarioParameters{
			DevicesDown: devicesDown,
			LinksDown:   linksDown,
		},
	}

	for _, device := range devicesDown {
		if _, ok := graph.Nodes[device]; !ok {
			return nil, fmt.Errorf("device not found: %s", device)
		}
		graph.Nodes[device].Status = false
	}

	for _, link := range linksDown {
		if _, ok := graph.Links[link]; !ok {
			return nil, fmt.Errorf("link not found: %s", link)
		}
		graph.Links[link].Status = false
	}

	defer graph.FullReset()
	graph.ComputeAllLinkStatus()

	isolated, err := graph.GetIsolatedBottomNodes()
	if err != nil {
		return nil, err
	}
	result.ImpactedNodes = isolated

	return &result, nil
}

func RunWithAssetsDown(graph *topology.Graph, devicesDown []string, linksDown []string) (*SimulationResult, error) {
	results := SimulationResult{
		ScenarioResults: make(map[string]*ScenarioResult),
	}

	defer utils.Timer(&results.Performance)()
	result, err := runScenario(graph, devicesDown, nil)
	if err != nil {
		return nil, err
	}

	results.ScenarioResults["custom"] = result

	return &results, nil
}

func RunAllNodesScenarios(graph *topology.Graph) (*SimulationResult, error) {
	results := SimulationResult{
		ScenarioResults: make(map[string]*ScenarioResult),
	}

	defer utils.Timer(&results.Performance)()
	for _, node := range graph.Nodes {
		if node.Role == configs.BottomDeviceRole {
			continue
		}
		result, err := runScenario(graph, []string{node.Hostname}, nil)
		if err != nil {
			return nil, err
		}

		results.ScenarioResults[node.Hostname] = result
	}

	return &results, nil
}
