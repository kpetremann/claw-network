package topology

import "testing"

func TestNodeUnmarshal(t *testing.T) {
	var node Node
	node.UnmarshalJSON([]byte(`{"hostname":"tor-01-01","role":"tor","layer":1,"status":true,"realStatus":true}`))

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

func TestNodeIsIsolated(t *testing.T) {
	node := Node{Hostname: "tor-01-01", Role: "tor", Layer: 1, Status: false, RealStatus: true}
	node.Uplinks = map[string]*Link{
		"spine1": {CanReachEdge: false},
		"spine2": {CanReachEdge: false},
	}

	if isolated, _ := node.IsIsolated(); isolated != true {
		t.Error("Did not detect isolated node")
	}

	node.Uplinks["spine1"].CanReachEdge = true
	if isolated, _ := node.IsIsolated(); isolated != false {
		t.Error("Node should not be considered isolated")
	}
}
