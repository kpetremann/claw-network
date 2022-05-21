package topology

import (
	"testing"
)

func TestNodeUnmarshal(t *testing.T) {
	var node Node
	if err := node.UnmarshalJSON([]byte(`{"hostname":"tor-01-01","role":"tor","layer":1,"status":true,"realStatus":true}`)); err != nil {
		t.Errorf("Test issues: unable to unmarshal the sample. %s", err)
	}

	if node.Hostname != "tor-01-01" {
		t.Error("Unmarshal failure: bad hostname")
	}

	if node.Layer != 1 {
		t.Error("Unmarshal failure: bad layer")
	}

	if node.Role != "tor" {
		t.Error("Unmarshal failure: bad role")
	}

	if node.Status != true {
		t.Error("Unmarshal failure: bad status")
	}

	if node.RealStatus != true {
		t.Error("Unmarshal failure: bad realStatus")
	}
}

func TestNodeResetStatus(t *testing.T) {
	node := Node{Hostname: "tor-01-01", Role: "tor", Layer: 1, Status: false, RealStatus: true}
	node.ResetStatus()

	if node.Status != true {
		t.Error("ResetStatus failure")
	}
}

func TestNodeComputeAllLinkStatus(t *testing.T) {
	graph := GenerateMinimumGraph()
	graph.Nodes["tor1"].ComputeAllLinkStatus(1)

	for _, uplink := range graph.Nodes["tor1"].Uplinks {
		if uplink.BuildId != 1 {
			t.Error("Not all uplinks have the same build id")
		}
	}
}

func TestNodeIsIsolatedFromTop(t *testing.T) {
	// ToR
	tor := Node{Hostname: "tor-01-01", Role: "tor", Layer: 1, Status: false, RealStatus: true}
	tor.Uplinks = map[string]*Link{
		"spine1": {CanReachEdge: false},
		"spine2": {CanReachEdge: false},
	}

	if isolated, _ := tor.IsIsolatedFromTop(); isolated == false {
		t.Error("Did not detect isolated node")
	}

	tor.Uplinks["spine1"].CanReachEdge = true
	if isolated, _ := tor.IsIsolatedFromTop(); isolated == true {
		t.Error("Node should not be considered isolated")
	}

	// Edge
	edge := Node{Hostname: "edge1", Role: "edge", Layer: 3, Status: false, RealStatus: true, Uplinks: map[string]*Link{}}
	if isolated, _ := edge.IsIsolatedFromTop(); isolated == true {
		t.Error("Edge should never be considered isolated as they are at the top of the graph")
	}
}

func TestNodeIsConnectedTor(t *testing.T) {
	tor := Node{Hostname: "tor-01-01", Role: "tor", Layer: 1, Status: false, RealStatus: true, Uplinks: map[string]*Link{}, Downlinks: map[string]*Link{}}
	if used, _ := tor.IsConnected(); used == true {
		t.Error("Bottom nodes without uplinks should not be considered connected")
	}

	tor.Uplinks["spine1"] = &Link{}
	if used, _ := tor.IsConnected(); used == false {
		t.Error("Bottom nodes without at least one uplink should be considered connected")
	}
}

func TestNodeIsConnectedSpine(t *testing.T) {
	spine := Node{Hostname: "spine1", Role: "spine", Layer: 2, Status: false, RealStatus: true, Uplinks: map[string]*Link{}, Downlinks: map[string]*Link{}}

	// 0 links
	if used, _ := spine.IsConnected(); used == true {
		t.Error("Spine nodes without links should not be considered connected")
	}

	// 1 uplink, 0 downlink
	spine.Uplinks["edge1"] = &Link{}
	if used, _ := spine.IsConnected(); used == true {
		t.Error("Spine nodes with without downlinks should not be considered connected")
	}

	// 1 uplink, 1 downlink
	spine.Downlinks["tor1"] = &Link{}
	if used, _ := spine.IsConnected(); used == false {
		t.Error("Spine nodes with both uplinks and downlinks should be considered connected")
	}

	// 0 uplink, 1 downlink
	delete(spine.Uplinks, "edge1")
	if used, _ := spine.IsConnected(); used == true {
		t.Error("Spine nodes with without uplinks should not be considered connected")
	}
}

func TestNodeIsConnectedEdge(t *testing.T) {
	edge := Node{Hostname: "edge1", Role: "edge", Layer: 3, Status: false, RealStatus: true, Uplinks: map[string]*Link{}, Downlinks: map[string]*Link{}}
	if used, _ := edge.IsConnected(); used == true {
		t.Error("Edges without downlinks should not be considered connected")
	}

	edge.Downlinks["spine1"] = &Link{}
	if used, _ := edge.IsConnected(); used == false {
		t.Error("Edges with at least one downlink should be considered connected")
	}
}
