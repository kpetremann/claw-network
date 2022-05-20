package topology

type Anomaly struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Anomalies struct {
	Node      string    `json:"node"`
	Anomalies []Anomaly `json:"anomalies"`
}
