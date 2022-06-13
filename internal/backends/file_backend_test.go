package backends

import "testing"

func TestGetTopologies(t *testing.T) {
	repository := &FileRepository{}
	repository.Topologies = []string{"topology1", "topology2"}
	topologies := repository.GetTopologies()

	if topologies[0] != "topology1" {
		t.Errorf("Expected topology1, got %s", topologies[0])
	}

	if topologies[1] != "topology2" {
		t.Errorf("Expected topology2, got %s", topologies[0])
	}
}
