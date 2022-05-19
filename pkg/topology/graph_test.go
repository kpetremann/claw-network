package topology

import (
	"encoding/json"
	"os"
	"testing"
)

func TestGraphUnmarshalJSON(t *testing.T) {
	expectedNodesNumber := 9
	expectedLinksNumber := 10

	expectedNodes := map[string]Node{
		"tor-01-01":   {Hostname: "tor-01-01", Role: "tor", Layer: 1, Status: true, RealStatus: true},
		"tor-01-02":   {Hostname: "tor-01-02", Role: "tor", Layer: 1, Status: true, RealStatus: true},
		"fabric-1-01": {Hostname: "fabric-1-01", Role: "fabric", Layer: 2, Status: true, RealStatus: true},
		"fabric-2-01": {Hostname: "fabric-2-01", Role: "fabric", Layer: 2, Status: true, RealStatus: true},
		"fabric-1-02": {Hostname: "fabric-1-02", Role: "fabric", Layer: 2, Status: true, RealStatus: true},
		"fabric-2-02": {Hostname: "fabric-2-02", Role: "fabric", Layer: 2, Status: true, RealStatus: true},
		"spine-11":    {Hostname: "spine-11", Role: "spine", Layer: 3, Status: true, RealStatus: true},
		"spine-12":    {Hostname: "spine-12", Role: "spine", Layer: 3, Status: true, RealStatus: true},
		"edge-1":      {Hostname: "edge-1", Role: "edge", Layer: 4, Status: true, RealStatus: true},
	}

	jsonData, err := os.ReadFile("../../examples/small_topology.json")
	if err != nil {
		t.Error("Test issue: topology sample not found")
		return
	}

	var graph Graph
	if err := json.Unmarshal(jsonData, &graph); err != nil {
		t.Error("Unmarshal failure")
		return
	}

	countUplinks := 0
	for uid, node := range graph.Nodes {
		expectedNode := expectedNodes[node.Hostname]

		// Check nodes
		if uid != node.Hostname {
			t.Errorf("Node %s has a wrong hostname: %s", uid, node.Hostname)
		}

		if node.Role != expectedNode.Role {
			t.Errorf("Node %s has wrong 'Role': %s != %s", node.Hostname, node.Role, expectedNode.Role)
		}
		if node.Layer != expectedNode.Layer {
			t.Errorf("Node %s has wrong 'Layer': %d != %d", node.Hostname, node.Layer, expectedNode.Layer)
		}
		if node.Status != expectedNode.Status {
			t.Errorf("Node %s has wrong 'Status': %t != %t", node.Hostname, node.Status, expectedNode.Status)
		}
		if node.RealStatus != expectedNode.RealStatus {
			t.Errorf("Node %s has wrong 'RealStatus': %t != %t", node.Hostname, node.RealStatus, expectedNode.RealStatus)
		}

		// Check uplinks of all nodes
		for _, uplink := range node.Uplinks {
			countUplinks++
			if node.Hostname != uplink.SouthNode.Hostname {
				t.Errorf("Link not attached to the right node: %s != %s", node.Hostname, uplink.SouthNode.Hostname)
			}
		}
	}

	// Count elements
	if len(graph.Nodes) != expectedNodesNumber {
		t.Error("Missing nodes")
	}

	if len(graph.Links) != expectedLinksNumber {
		t.Error("Missing links")
	}
	if countUplinks != expectedLinksNumber {
		t.Error("Some links were not attached to a node")
	}
}

func TestGraphMarshalJSON(t *testing.T) {
	expectedJSON := "{\"BuildId\":0,\"Links\":{\"1\":{\"uid\":\"1\",\"south_node\":\"tor1\",\"north_node\":\"tor1\",\"status\":false},\"2\":{\"uid\":\"2\",\"south_node\":\"tor1\",\"north_node\":\"tor1\",\"status\":false}},\"Nodes\":{\"spine1\":{\"hostname\":\"tor1\",\"layer\":2,\"role\":\"spine\",\"status\":false},\"tor1\":{\"hostname\":\"tor1\",\"layer\":1,\"role\":\"tor\",\"status\":false},\"tor2\":{\"hostname\":\"tor1\",\"layer\":1,\"role\":\"tor\",\"status\":false}},\"BottomNode\":null}"
	nodes := map[string]*Node{
		"tor1":   {Hostname: "tor1", Layer: 1, Role: "tor", Uplinks: map[string]*Link{}},
		"tor2":   {Hostname: "tor1", Layer: 1, Role: "tor", Uplinks: map[string]*Link{}},
		"spine1": {Hostname: "tor1", Layer: 2, Role: "spine", Uplinks: map[string]*Link{}},
	}
	links := map[string]*Link{
		"1": {Uid: "1", SouthNode: nodes["tor1"], NorthNode: nodes["spine1"]},
		"2": {Uid: "2", SouthNode: nodes["tor2"], NorthNode: nodes["spine1"]},
	}

	graph := Graph{
		Nodes: nodes,
		Links: links,
	}

	graphBytes, err := json.Marshal(graph)
	if err != nil {
		t.Errorf("Error while trying to marshal the Graph: %s", err)
	}

	if string(graphBytes) != expectedJSON {
		t.Errorf("Bad JSON result: %s", graphBytes)
	}
}

func TestGraphMarshalAndUnmarshal(t *testing.T) {
	graph := GenerateMinimumGraph()
	graphBytes, err := json.Marshal(graph)
	if err != nil {
		t.Errorf("Error while trying to marshal the Graph: %s", err)
	}

	var newGraph Graph
	if err := json.Unmarshal(graphBytes, &newGraph); err != nil {
		t.Errorf("Failed to reload generated JSON: %s", graphBytes)
	}
}

func TestGraphAddLink(t *testing.T) {
	graph := GenerateMinimumGraph()
	graph.Links = make(map[string]*Link) // removing links for the tests

	// Horizontal links are not supported
	if err := graph.AddLink(&LinkDefinition{Uid: "1", SouthNodeHostname: "tor1", NorthNodeHostname: "tor2"}); err == nil {
		t.Error("Horizontal links should not be permitted")
	}

	// Downlinks are not supported
	if err := graph.AddLink(&LinkDefinition{Uid: "1", SouthNodeHostname: "spine1", NorthNodeHostname: "tor1"}); err == nil {
		t.Error("Horizontal links should not be permitted")
	}

	// Uplinks are supported
	if err := graph.AddLink(&LinkDefinition{Uid: "1", SouthNodeHostname: "tor1", NorthNodeHostname: "spine1"}); err != nil {
		t.Error("Unable to add an uplink")
	}
}

func TestGraphComputeAllLinkStatus(t *testing.T) {
	graph := GenerateMinimumGraph()
	graph.ConnectNodesToLinks()

	graph.ComputeAllLinkStatus()

	if graph.BuildId != 1 {
		t.Error("Bad BuildID")
	}

	for _, link := range graph.Links {
		if link.BuildId != 1 {
			t.Error("Not all nodes has been computed")
			break
		}
	}
}

func TestGraphFullReset(t *testing.T) {
	graph := GenerateMinimumGraph()
	graph.Nodes["tor1"].Status = false
	graph.Links["1"].Status = false

	graph.FullReset()

	if graph.Nodes["tor1"].Status == false && graph.Nodes["tor1"].Status != graph.Nodes["tor1"].RealStatus {
		t.Error("Failed to reset node status")
	}

	if graph.Links["1"].Status == false && graph.Links["1"].Status != graph.Links["1"].RealStatus {
		t.Error("Failed to reset link status")
	}
}
