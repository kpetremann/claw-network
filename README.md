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

ClawNetwork can be leveraged to detect SPOF of any anomalies such as spine without downlinks.

# Quickstart

## From source

Simply run ClawNetwork app using `go run .`

Alternative: build the binary via `go build` and run it.

## Using Docker compose

### Default backend

Run ClawNetwork with default backend (FileRepository):
```shell
docker-compose -f compose/docker-compose.yml up -d
```

FileRepository stores the topologies in dedicated JSON files on the disk.

By default, this uses `examples/` directory provided in this repository.

> At the moment this is not customizable, but it will be very soon.

### Run with the Backend of your choice

```shell
docker-compose -f compose/docker-compose.yml -f <backend>.yml up -d
```

#### RedisJSON

> recommended backend for production if you need to store topologies

At the moment, Redis JSON is the only alternative backend:
```shell
docker-compose -f compose/docker-compose.yml -f redisjson.yml up -d
```

This backend leverages [RedisJSON module](https://redis.io/docs/stack/json/) to store pure JSON to Redis. Persistence is enabled and forced at each changes (ADD/DELETE) by ClawNetwork.

## Configuration

Configuration can be configured either via environment variables or YAML file (settings.yaml).

List of parameters available (`varenv format` | `YAML format`):

- `CLAW_LISTENADDRESS` | `ListenAddress`: ClawNetwork API listen address (default: `"0.0.0.0"`)
- `CLAW_LISTENPORT` | `ListenPort`: ClawNetwork API listen port (default: `"8080"`)
- `CLAW_TOPDEVICEROLE` | `TopDeviceRole`: Role of device at the top of the topology graph (default: `"edge"`)
- `CLAW_BOTTOMDEVICEROLE` | `BottomDeviceRole`: Role of device at the Bottom of the topology graph (default: `"tor"`)
- `CLAW_BACKEND` | `Backend`: Choose backend to store topologies (choices: `"file", "redis"`, default: `"file"`)
- `CLAW_BACKENDS.FILE.PATH` | `Backends.Redis.Path`: Redis DB to use (default: `"./topologies/"`)
- `CLAW_BACKENDS.REDIS.HOST` | `Backends.Redis.Host`: Redis server address (default: `"localhost"`)
- `CLAW_BACKENDS.REDIS.PORT` | `Backends.Redis.Port`: Redis server port (default: `"6379"`)
- `CLAW_BACKENDS.REDIS.PASSWORD` | `Backends.Redis.Password`: Redis password (default: `""`)
- `CLAW_BACKENDS.REDIS.DB` | `Backends.Redis.DB`: Redis DB to use (default: `0`)

# Usage

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

### Anomaly detection

- GET `/topology/:topology_name/anomalies`: get topology anomalies

It list all anomalies in the topology graph.

#### Link anomalies

A node is not connected properly to the graph.

For example:
- a ToR does not have any uplinks
- a spine does not have any downlinks or any uplinks
- an edge does not have any downlinks

This does not consider the status of the link, it only checks if there is a link.

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

# Integrations

Below some ideas of possible integrations:

- the client push the topology with the simulation request. The topology is not stored.
```
+-------------------------+
|  Observability metrics  |
|   example: Prometheus   |
+-------------------------+
             ^
             |
             | get metrics
             |
             |
             |
 +-----------------------+
 |                       |           get impact
 |        Client         |        on custom topology        +---------------+
 |   => convert metrics  |--------------------------------->|  ClawNetwork  |
 |      to topology      |                                  +---------------+
 +-----------------------+
```

- the client provides the topologies and they are stored
```
+-------------------------+
|  Observability metrics  |
|   example: Prometheus   |
+-------------------------+
             ^
             |
             | get metrics
             |
             |
             |
 +-----------------------+
 |        Client         |       push topology      +---------------+      save topology       +-------------------------+
 |   => convert metrics  |------------------------->|  ClawNetwork  |<------------------------>| Storage (FS, redis,...) |
 |      to topology      |        get impact        +---------------+       get topology       +-------------------------+
 +-----------------------+
```

- dedicated topology provider
```
                                                 +---------------------+
+-------------------------+                      |  Topology provider  |
|  Observability metrics  | <------------------- | => convert metrics  |
+-------------------------+                      |    to topology      |
                                                 +---------------------+
                                                            |
                                                            |
                                                            | push topology
                                                            |
                                                            |
                                                            |
                                                            v
 +-----------------------+        get impact        +---------------+      save topology       +-------------------------+
 |        Client         |------------------------->|  ClawNetwork  |<------------------------>| Storage (FS, redis,...) |
 +-----------------------+                          +---------------+       get topology       +-------------------------+
```


# Todo / coming features

See [Project board](https://github.com/kpetremann/claw-network/projects/1)
