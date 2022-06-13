package backends

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	. "github.com/kpetremann/claw-network/configs"
	"github.com/kpetremann/claw-network/pkg/topology"
)

const jsonSuffix = ".json"

type FileRepository struct {
	Topologies []string
	lock       sync.RWMutex
}

// Getter for Topologies
func (r *FileRepository) GetTopologies() []string {
	r.lock.RLock()
	defer r.lock.RUnlock()

	topologiesCopy := make([]string, len(r.Topologies))
	copy(topologiesCopy, r.Topologies)
	return topologiesCopy
}

// Get topology list from files
func (t *FileRepository) RefreshRepository() error {
	if _, err := os.Stat(Config.Backends.File.Path); os.IsNotExist(err) {
		if err := os.Mkdir(Config.Backends.File.Path, os.ModePerm); err != nil {
			panic(err)
		}
	}

	t.lock.Lock()
	defer t.lock.Unlock()

	files, err := os.ReadDir(Config.Backends.File.Path)
	if err != nil {
		return err
	}

	var newTopologyList []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if filepath.Ext(fileName) != jsonSuffix {
			continue
		}

		newTopologyList = append(newTopologyList, strings.TrimSuffix(fileName, jsonSuffix))
	}

	t.Topologies = newTopologyList

	return nil
}

func (t *FileRepository) LoadTopology(topologyFile string) (*topology.Graph, error) {
	var topo topology.Graph

	t.lock.RLock()
	byteValue, err := os.ReadFile(Config.Backends.File.Path + topologyFile + jsonSuffix)
	if err != nil {
		t.lock.RUnlock()
		return nil, err
	}
	t.lock.RUnlock()

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

	t.lock.Lock()
	if err := os.WriteFile(Config.Backends.File.Path+fileName+jsonSuffix, jsonTopology, 0644); err != nil {
		t.lock.Unlock()
		return err
	}
	t.lock.Unlock()

	// Refresh topology list
	if err := t.RefreshRepository(); err != nil {
		return err
	}

	return nil
}

func (t *FileRepository) DeleteTopology(topologyName string) error {
	t.lock.Lock()
	if err := os.Remove(Config.Backends.File.Path + topologyName + jsonSuffix); err != nil {
		t.lock.Unlock()
		return err
	}
	t.lock.Unlock()

	// Refresh topology list
	if err := t.RefreshRepository(); err != nil {
		return err
	}

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
