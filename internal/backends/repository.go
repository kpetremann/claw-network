package backends

import "github.com/kpetremann/claw-network/pkg/topology"

type Repository interface {
	GetTopologies() []string
	SetTopologies(topologies []string)
	RefreshRepository() error
	LoadTopology(topologyName string) (*topology.Graph, error)
	SaveTopology(name string, graph *topology.Graph) error
	DeleteTopology(topologyName string) error
	GetTopologyDetails(topologyName string) (map[string]int, error)
	ListTopologiesDetails() (map[string]map[string]int, error)
}

func MigrateRepository(currentRepo, newRepo Repository) error {
	if err := currentRepo.RefreshRepository(); err != nil {
		return err
	}

	for _, topologyName := range currentRepo.GetTopologies() {
		topology, err := currentRepo.LoadTopology(topologyName)
		if err != nil {
			return err
		}

		if err := newRepo.SaveTopology(topologyName, topology); err != nil {
			return err
		}
	}

	return nil
}
