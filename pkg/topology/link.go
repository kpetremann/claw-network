package topology

import (
	"encoding/json"
	"errors"

	. "github.com/kpetremann/claw-network/configs"
)

type Link struct {
	Uid              string
	SouthNode        *Node
	NorthNode        *Node
	Status           bool
	RealStatus       bool
	CanReachEdge     bool
	RealCanReachEdge bool
	BuildId          int
}

type LinkDefinition struct {
	Uid               string `json:"uid"`
	SouthNodeHostname string `json:"south_node"`
	NorthNodeHostname string `json:"north_node"`
	Status            bool   `json:"status"`
}

func (l *Link) UnmarshalJSON(_ []byte) error {
	return errors.New("cannot unmashal Link directly, must unmarshal via Graph")
}

func (l *Link) MarshalJSON() ([]byte, error) {
	jLink := LinkDefinition{
		Uid:               l.Uid,
		SouthNodeHostname: l.SouthNode.Hostname,
		NorthNodeHostname: l.NorthNode.Hostname,
		Status:            l.Status,
	}
	return json.Marshal(&jLink)
}

func (l *Link) ResetStatus() {
	l.Status = l.RealStatus
	l.CanReachEdge = l.RealCanReachEdge
}

func (l *Link) ComputeNorthPathStatus(buildId int) bool {
	// if already built, we stop here
	if l.BuildId == buildId {
		return l.CanReachEdge
	}

	// if not uplink on the north device, we stop here
	if len(l.NorthNode.Uplinks) == 0 {
		l.BuildId = buildId

		// check if edge has been reached
		if l.NorthNode.Role == Config.TopDeviceRole {
			l.CanReachEdge = l.Status
		} else {
			l.CanReachEdge = false
		}

		return l.CanReachEdge
	}

	// check recursively all upper links
	// the current link is valid if at least one uplink is ok
	hasOneValidUplink := false
	for _, uplink := range l.NorthNode.Uplinks {
		hasOneValidUplink = hasOneValidUplink || uplink.ComputeNorthPathStatus(buildId)
	}

	// the link is valid only if the link is up, the nodes are up
	// and if there at least one valid uplink
	l.CanReachEdge = l.Status && hasOneValidUplink && l.NorthNode.Status && l.SouthNode.Status

	l.BuildId = buildId

	return l.CanReachEdge
}
