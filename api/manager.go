package api

import (
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
		for {
			select {
			case s.getRepository <- repository:
			case repository.Topologies = <-s.writeRepository:
			}
		}
	}()

	return s
}
