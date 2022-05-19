package topology

import (
	"encoding/json"
	"testing"
)

func TestLinkUnmarshalJSON(t *testing.T) {
	var link Link
	jsonLink := []byte(
		`{"uid":"tor1->spine1","southNode":{"hostname":"tor-01-01","role":"tor","layer":1,"status":true,"realStatus":true},"northNode":{"hostname":"spine-01-01","role":"spine","layer":2,"status":true,"realStatus":true},"status":true,"realStatus":true}`)
	if err := json.Unmarshal(jsonLink, &link); err == nil {
		t.Error("Links are not supposed to be unmarshaled directly")
	}
}

func TestLinkMarshalJSON(t *testing.T) {
	link := LinkDefinition{
		Uid:               "tor1->spine1",
		SouthNodeHostname: "tor-01-01",
		NorthNodeHostname: "spine-01-01",
		Status:            true,
	}

	expectedJSON := []byte(`{"uid":"tor1->spine1","south_node":"tor-01-01","north_node":"spine-01-01","status":true}`)
	jsonLink, err := json.Marshal(link)
	if err != nil {
		t.Errorf("Marshal failure: %s", err)
	}

	if string(jsonLink) != string(expectedJSON) {
		t.Errorf("Marshal failure: %s != %s", jsonLink, expectedJSON)
	}
}

func TestLinkResetStatus(t *testing.T) {
	graph := GenerateMinimumGraph()
	graph.Links["tor1->spine1"].Status = false

	graph.Links["tor1->spine1"].ResetStatus()

	if graph.Links["tor1->spine1"].Status == false {
		t.Error("Link status was not reset")
	}
}

func TestLinkComputeNorthPathStatus(t *testing.T) {
	graph := GenerateMinimumGraph()
	graph.ConnectNodesToLinks()

	// all should be ok
	graph.Links["tor1->spine1"].ComputeNorthPathStatus(1)

	if graph.Links["tor1->spine1"].BuildId != 1 {
		t.Error("Link was not computed")
	}

	if graph.Links["tor1->spine1"].CanReachEdge == false {
		t.Error("Link should be able to reach edge")
	}

	// removing path to edge
	graph.Links["spine1->edge1"].Status = false

	graph.Links["tor1->spine1"].ComputeNorthPathStatus(2)

	if graph.Links["tor1->spine1"].BuildId != 2 {
		t.Error("Link was not computed")
	}

	if graph.Links["tor1->spine1"].CanReachEdge == true {
		t.Error("Link should not be able to reach edge")
	}
}
