package main

import "fmt"

func testListTopology() {
	printTitle("Testing list topologies")
	var err error
	body, _, err := request("topology", "GET", nil)

	if err != nil {
		printStatus(fmt.Sprintf("error: %s", err))
		return
	}

	printJson(string(body))
	printStatus("Successfully listed topologies")
}

func testGetTopology() {
	printTitle("Testing get topology")
	var err error
	body, _, err := request("topology/small_topology", "GET", nil)

	if err != nil {
		printStatus(fmt.Sprintf("error: %s", err))
		return
	}

	printJson(string(body))
	printStatus("Successfully got topology")
}

func testAddAndDeleteTopology() {
	printTitle("Testing add topology")

	topo := []byte(`{"nodes":[{"hostname":"tor-01-01","role":"tor","status":true,"layer":1},{"hostname":"fabric-1-01","role":"fabric","status":true,"layer":2}],"links":[{"south_node":"tor-01-01","north_node":"fabric-1-01","status":true,"uid":"10.0.0.0->10.0.0.1"}]}`)

	var err error
	body, _, err := request("topology/test", "POST", topo)

	if err != nil {
		printStatus(fmt.Sprintf("error: %s, details: %s", err, string(body)))
		return
	}

	printJson(string(body))
	printStatus("Successfully added topology")

	printTitle("Testing delete topology")
	body, _, err = request("topology/test", "DELETE", nil)

	if err != nil {
		printStatus(fmt.Sprintf("error: %s", err))
		return
	}

	printJson(string(body))
	printStatus("Successfully deleted topology")
}
