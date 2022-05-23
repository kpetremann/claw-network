// nolint
package main

import (
	"fmt"

	"github.com/kpetremann/claw-network/internal/backends"
)

// Get all topologies from FileBackend repository
// Save them to RedisRepository
func MigrateFileToRedis() {
	fileRepository := backends.TopologyRepository{}
	redisRepository := backends.RedisRepository{}

	if err := backends.MigrateRepository(&fileRepository, &redisRepository); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Migration complete")
	}
}

func main() {
	MigrateFileToRedis()
	redisRepository := backends.RedisRepository{}

	redisRepository.RefreshRepository()
	d, _ := redisRepository.ListTopologiesDetail()
	fmt.Println(d)

	t, _ := redisRepository.LoadTopology("small_topology")
	fmt.Println(t)

	redisRepository.DeleteTopology("small_topology")
}
