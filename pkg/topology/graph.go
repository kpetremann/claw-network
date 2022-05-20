package topology

import (
	"encoding/json"
	"fmt"

	"github.com/kpetremann/claw-network/configs"
)

type LinkNotSupportedError struct {
	Link    *LinkDefinition
	Message string
}

func (e *LinkNotSupportedError) Error() string {
	return "Unable to create link " + e.Link.Uid + " because " + e.Message
}

type Graph struct {
	BuildId    int
	Links      map[string]*Link
	Nodes      map[string]*Node
	BottomNode []*Node
}

func (g *Graph) String() string {
	repr := "nodes: "
	for _, node := range g.Nodes {
		repr += fmt.Sprintf("%v ", node)
	}
	return repr
}

func (g *Graph) UnmarshalJSON(data []byte) error {
	jGraph := struct {
		Nodes []*Node           `json:"nodes"`
		Links []*LinkDefinition `json:"links"`
	}{}

	if err := json.Unmarshal(data, &jGraph); err != nil {
		return err
	}

	// extract nodes
	g.Nodes = make(map[string]*Node)
	for _, node := range jGraph.Nodes {
		g.Nodes[node.Hostname] = node
		if node.Role == configs.BottomDeviceRole {
			g.BottomNode = append(g.BottomNode, node)
		}
	}

	// create links
	g.Links = make(map[string]*Link)
	for _, linkDef := range jGraph.Links {
		if err := g.AddLink(linkDef); err != nil {
			return err
		}
	}

	return nil
}

func (g *Graph) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Nodes []*Node `json:"nodes"`
		Links []*Link `json:"links"`
	}{
		Nodes: g.GetNodeList(),
		Links: g.GetLinkList(),
	})
}

func (g *Graph) AddLink(linkDef *LinkDefinition) error {
	southNode := g.Nodes[linkDef.SouthNodeHostname]
	northNode := g.Nodes[linkDef.NorthNodeHostname]

	if valid, reason := g.IsLinkSupported(southNode, northNode); !valid {
		return &LinkNotSupportedError{linkDef, reason}
	}

	link := Link{
		SouthNode:  southNode,
		NorthNode:  northNode,
		Uid:        linkDef.Uid,
		Status:     linkDef.Status,
		RealStatus: linkDef.Status,
	}

	// connect node to links
	link.SouthNode.Uplinks[linkDef.Uid] = &link
	link.NorthNode.Downlinks[linkDef.Uid] = &link

	g.Links[link.Uid] = &link

	return nil
}

func (g *Graph) ConnectNodesToLinks() {
	for _, link := range g.Links {
		link.SouthNode.Uplinks[link.Uid] = link
		link.NorthNode.Downlinks[link.Uid] = link
	}
}

func (g *Graph) IsLinkSupported(southNode *Node, northNode *Node) (bool, string) {
	if northNode == nil || southNode == nil {
		return false, "referring to unknown node"
	}
	if northNode.Layer < southNode.Layer {
		return false, "cyclic graph links are not supported"
	}
	if northNode.Layer != southNode.Layer+1 {
		return false, "horizontal links are not supported"
	}

	return true, ""
}

func (g *Graph) GetNodeList() []*Node {
	allNodes := make([]*Node, 0, len(g.Nodes))

	for _, node := range g.Nodes {
		allNodes = append(allNodes, node)
	}

	return allNodes
}

func (g *Graph) GetLinkList() []*Link {
	allLinks := make([]*Link, 0, len(g.Links))

	for _, link := range g.Links {
		allLinks = append(allLinks, link)
	}

	return allLinks
}

func (g *Graph) ComputeAllLinkStatus() {
	g.BuildId += 1

	for _, tor := range g.BottomNode {
		tor.ComputeAllLinkStatus(g.BuildId)
	}
}

func (g *Graph) GetIsolatedBottomNodes() ([]string, error) {
	var BottomNode []string
	for _, tor := range g.BottomNode {
		isIsolated, err := tor.IsIsolatedFromTop()
		if err != nil {
			return nil, err
		}
		if isIsolated {
			BottomNode = append(BottomNode, tor.Hostname)
		}
	}
	return BottomNode, nil
}

func (g *Graph) FullReset() {
	for _, link := range g.Links {
		link.ResetStatus()
	}

	for _, node := range g.Nodes {
		node.ResetStatus()
	}
}
