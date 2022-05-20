package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/charmbracelet/glamour"
)

const BASE_URL = "http://localhost:8080/"

var output string = ""

func printTitle(in string) {
	toRender := fmt.Sprintf("# %s", in)
	output += "\n" + toRender
}

func printJson(in string) {
	toRender := fmt.Sprintf("```json\n%s\n```", in)
	output += "\n" + toRender

}

func printStatus(in string) {
	toRender := fmt.Sprintf("> %s", in)
	output += "\n" + toRender
}

func request(endpoint string, method string, postData []byte) ([]byte, interface{}, error) {
	var resp *http.Response
	var err error

	switch method {
	case "GET":
		resp, err = http.Get(BASE_URL + endpoint)
	case "POST":
		resp, err = http.Post(BASE_URL+endpoint, "application/json", bytes.NewReader(postData))
	case "DELETE":
		req, _ := http.NewRequest("DELETE", BASE_URL+endpoint, nil)
		resp, err = http.DefaultClient.Do(req)
	}

	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		return body, nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var parsed interface{}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return body, nil, err
	}

	return body, parsed, nil
}

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

func main() {
	testListTopology()
	testGetTopology()
	testAddAndDeleteTopology()

	out, _ := glamour.Render(output, "dark")
	fmt.Print(out)
}
