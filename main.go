package main

import (
	"github.com/kpetremann/claw-network/api"
	"github.com/kpetremann/claw-network/configs"
)

func main() {
	if err := configs.LoadConfig(); err != nil {
		panic(err)
	}
	api.ListenAndServer()
}
