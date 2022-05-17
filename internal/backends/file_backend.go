package backends

import (
	"encoding/json"
	"io/ioutil"

	"github.com/kpetremann/claw-network/configs"
	"github.com/kpetremann/claw-network/pkg/topology"
)

// Load topology information from JSON file
func LoadTopologyFromFile(topologyFile string) (*topology.Graph, error) {
	var topo topology.Graph
	byteValue, err := ioutil.ReadFile(configs.TopologyBaseDir + topologyFile)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byteValue, &topo); err != nil {
		return nil, err
	}

	return &topo, nil
}

// Save topology in a JSON file
func SaveTopologyToFile(fileName string, graph *topology.Graph) error {
	jsonTopology, err := json.Marshal(graph)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(configs.TopologyBaseDir+fileName, jsonTopology, 0644); err != nil {
		return err
	}

	return nil
}
