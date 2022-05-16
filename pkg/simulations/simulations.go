package simulations

import "github.com/kpetremann/claw-network/pkg/topology"

func AssetsDown(graph *topology.Graph, devicesDown []string, linksDown []string) ([]string, error) {
	for _, device := range devicesDown {
		graph.Nodes[device].Status = false
	}

	for _, link := range linksDown {
		graph.Links[link].Status = false
	}

	defer graph.FullReset()
	graph.ComputeAllLinkStatus()

	isolated, err := graph.GetIsolatedBottomNode()
	if err != nil {
		return nil, err
	}

	return isolated, nil
}
