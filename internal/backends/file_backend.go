package backends

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/kpetremann/claw-network/configs"
	"github.com/kpetremann/claw-network/pkg/topology"
)

// Load topology information from JSON file
func LoadTopologyFromFile(topologyFile string) (*topology.Graph, error) {
	var topo topology.Graph
	jsonFile, err := os.Open(configs.TopologyBaseDir + topologyFile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &topo)
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
