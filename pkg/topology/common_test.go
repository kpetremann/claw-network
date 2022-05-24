package topology

import (
	. "github.com/kpetremann/claw-network/configs"
)

func GenerateMinimumGraph() *Graph {
	LoadTestConfig()
	nodes := map[string]*Node{
		"tor1":   {Hostname: "tor1", Layer: 1, Role: "tor", Uplinks: map[string]*Link{}, Downlinks: map[string]*Link{}, Status: true, RealStatus: true},
		"tor2":   {Hostname: "tor2", Layer: 1, Role: "tor", Uplinks: map[string]*Link{}, Downlinks: map[string]*Link{}, Status: true, RealStatus: true},
		"spine1": {Hostname: "spine1", Layer: 2, Role: "spine", Uplinks: map[string]*Link{}, Downlinks: map[string]*Link{}, Status: true, RealStatus: true},
		"edge1":  {Hostname: "edge1", Layer: 3, Role: "edge", Uplinks: map[string]*Link{}, Downlinks: map[string]*Link{}, Status: true, RealStatus: true},
	}
	links := map[string]*Link{
		"tor1->spine1":  {Uid: "1", SouthNode: nodes["tor1"], NorthNode: nodes["spine1"], Status: true, RealStatus: true},
		"tor2->spine1":  {Uid: "2", SouthNode: nodes["tor2"], NorthNode: nodes["spine1"], Status: true, RealStatus: true},
		"spine1->edge1": {Uid: "3", SouthNode: nodes["spine1"], NorthNode: nodes["edge1"], Status: true, RealStatus: true},
	}

	graph := Graph{
		Nodes:      nodes,
		Links:      links,
		BottomNode: []*Node{nodes["tor1"], nodes["tor2"]},
	}

	return &graph
}

func LoadTestConfig() {
	Config.ListenAddress = "127.0.0.1"
	Config.ListenPort = "8080"
	Config.TopDeviceRole = "edge"
	Config.BottomDeviceRole = "tor"
	Config.Backend = "File"
	Config.Backends.File.Path = "./examples/"
}
