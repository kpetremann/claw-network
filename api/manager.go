package api

import (
	"fmt"
	"strings"

	. "github.com/kpetremann/claw-network/configs"
	"github.com/kpetremann/claw-network/internal/backends"
)

type SimulationManager struct {
	repository backends.Repository
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
	s := &SimulationManager{getRepository()}

	fmt.Println("Loading repository")
	if err := s.repository.RefreshRepository(); err != nil {
		panic(err)
	}

	return s
}
