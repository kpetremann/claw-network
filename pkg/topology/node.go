package topology

import (
	"encoding/json"
	"fmt"

	"github.com/kpetremann/claw-network/configs"
)

type Node struct {
	Hostname   string           `json:"hostname"`
	Uplinks    map[string]*Link `json:"-"`
	Downlinks  map[string]*Link `json:"-"`
	Layer      int              `json:"layer"`
	Role       string           `json:"role"`
	Status     bool             `json:"status"`
	RealStatus bool             `json:"-"`
}

func (n *Node) UnmarshalJSON(jsonText []byte) error {
	type NodeAlias Node
	var tmp NodeAlias

	if err := json.Unmarshal(jsonText, &tmp); err != nil {
		return err
	}

	*n = Node(tmp)

	n.Uplinks = make(map[string]*Link)
	n.Downlinks = make(map[string]*Link)
	n.RealStatus = n.Status

	return nil
}

func (n *Node) ResetStatus() {
	n.Status = n.RealStatus
}

func (n *Node) ComputeAllLinkStatus(buildId int) {
	for _, uplink := range n.Uplinks {
		uplink.ComputeNorthPathStatus(buildId)
	}
}

func (n *Node) IsIsolatedFromTop() (bool, error) {
	if n.Role == configs.TopDeviceRole {
		return false, nil
	}

	if len(n.Uplinks) == 0 {
		return true, fmt.Errorf("no uplink found on %s", n.Hostname)
	}

	// not isolated if at least one uplink is valid
	for _, uplink := range n.Uplinks {
		if uplink.CanReachEdge {
			return false, nil
		}
	}

	return true, nil
}

// Tells if a node is connected to the graph, without considering link status
// - a node is considered connected if it has at least one uplink and one downlink
// - a bottom node is considered connected if it has at least one uplink
// - a top node is considered connected if it has at least one downlink
func (n *Node) IsConnected() (bool, error) {
	if len(n.Uplinks) == 0 && len(n.Downlinks) == 0 {
		return false, fmt.Errorf("no link found on %s", n.Hostname)
	}
	if n.Role != configs.TopDeviceRole && len(n.Uplinks) == 0 {
		return false, fmt.Errorf("no uplink found on %s", n.Hostname)
	}
	if n.Role != configs.BottomDeviceRole && len(n.Downlinks) == 0 {
		return false, fmt.Errorf("no downlink found on %s", n.Hostname)
	}

	return true, nil
}
