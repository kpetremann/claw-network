package main

import "fmt"

func testAnomalies() {
	printTitle("Testing get topology anomalies")
	var err error
	body, _, err := request("topology/with_anomalies/anomalies", "GET", nil)

	if err != nil {
		printStatus(fmt.Sprintf("error: %s", err))
		return
	}

	printJson(string(body))
	printStatus("Successfully listed topology anomalies")
}
