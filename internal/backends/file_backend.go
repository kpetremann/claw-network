package backends

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/kpetremann/claw-network/configs"
	"github.com/kpetremann/claw-network/pkg/topology"
)

const jsonSuffix = ".json"

type TopologyRepository struct {
	Topologies []string
}

func (t *TopologyRepository) UpdateTopology() error {
	files, err := os.ReadDir(configs.TopologyBaseDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if filepath.Ext(fileName) != jsonSuffix {
			continue
		}

		t.Topologies = append(t.Topologies, strings.TrimSuffix(fileName, jsonSuffix))
	}

	return nil
}

func (t *TopologyRepository) DeleteTopology(topologyName string) error {
	if err := os.Remove(configs.TopologyBaseDir + topologyName + jsonSuffix); err != nil {
		return err
	}

	return nil
}

func (t *TopologyRepository) LoadTopology(topologyFile string) (*topology.Graph, error) {
	var topo topology.Graph
	// TODO: check if in repository before trying to read file
	byteValue, err := os.ReadFile(configs.TopologyBaseDir + topologyFile + jsonSuffix)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byteValue, &topo); err != nil {
		return nil, err
	}

	return &topo, nil
}

func (t *TopologyRepository) SaveTopology(fileName string, graph *topology.Graph) error {
	jsonTopology, err := json.Marshal(graph)
	if err != nil {
		return err
	}

	if err := os.WriteFile(configs.TopologyBaseDir+fileName, jsonTopology, 0644); err != nil {
		return err
	}

	return nil
}
