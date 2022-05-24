package api

import (
	"fmt"
	"strings"

	. "github.com/kpetremann/claw-network/configs"
	"github.com/kpetremann/claw-network/internal/backends"
)

type SimulationManager struct {
	getRepository   chan backends.Repository
	writeRepository chan []string
}

func getRepository() backends.Repository {
	switch strings.ToLower(Config.Backend) {
	case "file":
		fmt.Println("Using File backend")
		var repository backends.FileRepository
		return &repository
	case "redis":
		fmt.Println("Using Redis backend")
		var repository backends.RedisRepository
		return &repository
	default:
		panic("Unknown backend")
	}
}

func NewSimulationManager() *SimulationManager {
	s := &SimulationManager{
		getRepository:   make(chan backends.Repository),
		writeRepository: make(chan []string),
	}

	go func() {
		repository := getRepository()

		fmt.Println("Loading repository")
		if err := repository.RefreshRepository(); err != nil {
			panic(err)
		}

		for {
			select {
			case s.getRepository <- repository:
			case newTopologyList := <-s.writeRepository:
				repository.SetTopologies(newTopologyList)
			}
		}
	}()

	return s
}
