<img align="right" width="320px" src="https://raw.githubusercontent.com/kpetremann/claw-network/main/img/ClawNetwork-logo.png" />

[![status](https://img.shields.io/badge/status-in%20development-orange)](https://github.com/kpetremann/claw-network/)
[![Go](https://img.shields.io/github/go-mod/go-version/kpetremann/claw-network)](https://github.com/kpetremann/claw-network/)
[![GitHub](https://img.shields.io/github/license/kpetremann/claw-network)](https://github.com/kpetremann/claw-network/blob/main/LICENSE)

# Overview

ClawNetwork is a tool to simulate a network and evaluate failures impacts on Top of Racks.

It has been specially crafted for Clos Matrix network. For now, cyclic graphs are not supported. Only trees are.

```
Important notice:

This is in development and not fully usable yet.
But you can play with it :)
```

# Usecases


## Operations

The main usecase it to evaluate if an operation on a device in your core network will impact a Top of Rack.

Concerned operations can be: upgrade, reboot, risky maintenance etc...

## Detect anomalies / SPOF

ClawNetwork can be leveraged to detect SPOF of any anomalies.

Anomaly detection is an upcoming feature: detect if a node has no uplinks, or if they are not connected to anything...

# Usage

## Quickstart

Simply run ClawNetwork app using `go run .`

Alternative: build the binary via `go build` and run it.

### Manage stored topologies

- GET `/topology`: list stored topologies
- GET `/topology/:topology_name`: get topology definition
- POST `/topology/:topology_name`: create a new topology
- DELETE `/topology/:topology_name`: delete a topology

### Simulation on a stored topology

- GET `/topology/:topology_name/device/:device/down/impact`: run simulations on existing topology
- POST `/topology/custom/device/:device/down/impact`: run simulations on topology provided in the request body

It will run a simulation on a stored topology.

If `:device` is set to `each`, it will simulate failure impact of each devices excluding Top of Racks.

### Topology structure

The topology to provide looks like this in JSON:

```json
{
  "nodes": [
    {
      "hostname": "tor-01-01",
      "role": "tor",
      "status": true,
      "layer": 1
    },
    {
      "hostname": "fabric-1-01",
      "role": "fabric",
      "status": true,
      "layer": 2
    }
  ],
  "links": [
    {
      "south_node": "tor-01-01",
      "north_node": "fabric-1-01",
      "status": true,
      "uid": "10.0.0.0->10.0.0.1"
    }
  ]
}
```

> This structure is subject to change, as the API is not considered stable at the moment

### Example

Topology = 4 healthy fabric nodes + 4 healthy ToR

Simulations:
- first simulation considering first fabric node as down
- second simulation considering second fabric node as down but with the first up
- ...

# Example usecase

You can query the following endpoint to simulate down impact of each devices. It get the tppology example from the `example/full_topology_with_issues.json`.

```shell
$ curl http://127.0.0.1:8080/topology/full_topology_with_issues/device/each/down/impact | jq
{
  "scenarios_result": {
    "edge-0": {
      "impacts": null,
      "parameters": {
        "devices_down": [
          "edge-0"
        ],
        "links_down": null
      }
    },
    "edge-1": {
      "impacts": null,
      "parameters": {
        "devices_down": [
          "edge-1"
        ],
        "links_down": null
      }
    },
    "fabric-1-01": {
      "impacts": [
        "tor-01-01"
      ],
      "parameters": {
        "devices_down": [
          "fabric-1-01"
        ],
        "links_down": null
      }
    },
    ...,
    "compute_time": "89 ms"
}
```

As you can see, `tor-01-01` would be down if we shut `fabric-1-01`.

The topology defined in `example/full_topology_with_issues.json`, also specifies some devices as down. Here all the fabric of pod 01 has been set to down except for `fabric-1-01`. This is why if there is a failure on this device, it will impact `tor-01-01` as this ToR only had one healthy uplink.

Note: more advanced examples will be provided soon, with more complex scenarios.

# Todo / coming features

See [Project board](https://github.com/kpetremann/claw-network/projects/1)
