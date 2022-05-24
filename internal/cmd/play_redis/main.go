// nolint
package main

import (
	"fmt"
	"os"

	"github.com/kpetremann/claw-network/configs"
	"github.com/kpetremann/claw-network/internal/backends"
)

// Get all topologies from FileBackend repository
// Save them to RedisRepository
func MigrateFileToRedis() {
	fileRepository := backends.FileRepository{}
	redisRepository := backends.RedisRepository{}

	if err := backends.MigrateRepository(&fileRepository, &redisRepository); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Migration complete")
	}
}

func main() {
	os.Setenv("CLAW_BACKENDS.FILE.PATH", "./examples/")
	if err := configs.LoadConfig(); err != nil {
		panic(err)
	}
	MigrateFileToRedis()
	redisRepository := backends.RedisRepository{}

	redisRepository.RefreshRepository()
	d, _ := redisRepository.ListTopologiesDetails()
	fmt.Println(d)

	t, _ := redisRepository.LoadTopology("small_topology")
	fmt.Println(t)

	redisRepository.DeleteTopology("small_topology")
}
