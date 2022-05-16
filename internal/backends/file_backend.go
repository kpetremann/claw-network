package backends

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/kpetremann/claw-network/pkg/topology"
)

// Load topology information from JSON file
func LoadTopologyFromFile(topologyFile string) (*topology.Graph, error) {
	var topo topology.Graph
	jsonFile, err := os.Open(topologyFile)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &topo)
	return &topo, nil
}
