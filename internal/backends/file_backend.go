package backends

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	. "github.com/kpetremann/claw-network/configs"
	"github.com/kpetremann/claw-network/pkg/topology"
)

const jsonSuffix = ".json"

type FileRepository struct {
	Topologies []string
}

// Getter for Topologies
func (r *FileRepository) GetTopologies() []string {
	return r.Topologies
}

// Setter for Topologies
func (r *FileRepository) SetTopologies(topologies []string) {
	r.Topologies = topologies
}

// Get topology list from files
func (t *FileRepository) RefreshRepository() error {
	if _, err := os.Stat(Config.Backends.File.Path); os.IsNotExist(err) {
		if err := os.Mkdir(Config.Backends.File.Path, os.ModePerm); err != nil {
			panic(err)
		}
	}

	files, err := os.ReadDir(Config.Backends.File.Path)
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

func (t *FileRepository) LoadTopology(topologyFile string) (*topology.Graph, error) {
	var topo topology.Graph
	// TODO: check if in repository before trying to read file
	byteValue, err := os.ReadFile(Config.Backends.File.Path + topologyFile + jsonSuffix)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(byteValue, &topo); err != nil {
		return nil, err
	}

	return &topo, nil
}

func (t *FileRepository) SaveTopology(fileName string, graph *topology.Graph) error {
	jsonTopology, err := json.Marshal(graph)
	if err != nil {
		return err
	}

	if err := os.WriteFile(Config.Backends.File.Path+fileName+jsonSuffix, jsonTopology, 0644); err != nil {
		return err
	}
	t.Topologies = append(t.Topologies, fileName)

	return nil
}

func (t *FileRepository) DeleteTopology(topologyName string) error {
	if err := os.Remove(Config.Backends.File.Path + topologyName + jsonSuffix); err != nil {
		return err
	}

	// find element in the slice
	var index int
	for i, name := range t.Topologies {
		if topologyName == name {
			index = i
			break
		}
	}

	// delete the element
	t.Topologies = append(t.Topologies[:index], t.Topologies[index+1:]...)

	return nil
}

func (t *FileRepository) GetTopologyDetails(topologyName string) (map[string]int, error) {
	topology, err := t.LoadTopology(topologyName)
	if err != nil {
		return nil, err
	}

	results := make(map[string]int)
	results["nodes_total"] = len(topology.Nodes)
	results["links_total"] = len(topology.Links)

	for _, node := range topology.Nodes {
		if !node.Status {
			results["node_down"]++
		}
	}

	for _, link := range topology.Links {
		if !link.Status {
			results["link_down"]++
		}
	}

	results["nodes_up"] = results["nodes_total"] - results["node_down"]
	results["links_up"] = results["links_total"] - results["link_down"]

	return results, nil
}

func (t *FileRepository) ListTopologiesDetails() (map[string]map[string]int, error) {
	var err error
	topologies := make(map[string]map[string]int)

	for _, topologyName := range t.Topologies {
		topologies[topologyName], err = t.GetTopologyDetails(topologyName)
		if err != nil {
			return nil, err
		}
	}

	return topologies, nil
}
