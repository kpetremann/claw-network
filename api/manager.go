package api

import (
	"fmt"

	"github.com/kpetremann/claw-network/internal/backends"
)

type SimulationManager struct {
	getRepository   chan backends.TopologyRepository
	writeRepository chan []string
}

func NewSimulationManager() *SimulationManager {
	s := &SimulationManager{
		getRepository:   make(chan backends.TopologyRepository),
		writeRepository: make(chan []string),
	}

	go func() {
		var repository backends.TopologyRepository
		fmt.Println("Loading repository")
		if err := repository.UpdateTopology(); err != nil {
			panic(err)
		}

		for {
			select {
			case s.getRepository <- repository:
			case repository.Topologies = <-s.writeRepository:
			}
		}
	}()

	return s
}
