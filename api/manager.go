package api

import (
	"fmt"

	"github.com/kpetremann/claw-network/internal/backends"
)

type SimulationManager struct {
	getRepository   chan backends.FileRepository
	writeRepository chan []string
}

func NewSimulationManager() *SimulationManager {
	s := &SimulationManager{
		getRepository:   make(chan backends.FileRepository),
		writeRepository: make(chan []string),
	}

	go func() {
		var repository backends.FileRepository
		fmt.Println("Loading repository")
		if err := repository.RefreshRepository(); err != nil {
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
