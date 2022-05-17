# ClawNetwork

![status](https://img.shields.io/badge/status-in%20development-orange)
![Go](https://img.shields.io/github/go-mod/go-version/kpetremann/claw-network)
![GitHub](https://img.shields.io/github/license/kpetremann/claw-network)

## Overview

ClawNetwork is a tool to simulate a network and evaluate failures impacts on Top of Racks.

It has been specially crafted for Clos Matrix network.

```
Important notice:

This is in development and not fully usable yet. But you can play with it :)
```

## Usage example

First you just need to run the API, either using `go run .` or build it and run the executable.

You can query the following endpoint to simulate down impact of each devices. It get the tpology example from the `topology.json` file.

```
curl http://127.0.0.1:8080/topology/example/device/every/down/impact | jq
{
    "impact_simulation": {
        "edge-0": null,
        "edge-1": null,
        "fabric-1-01": [
            "tor-01-01"
        ],
        "fabric-1-02": null,
        ...
    }
}
```

As you can see, `tor-01-01` would be down if we shut `fabric-1-01`.

Note: more advanced examples will be provided soon, with more complex scenarios.

## Todo / coming features

- add tests
- get/list/push/delete a new topology
- simulate impact of link shutdown
- simulate impact of multiple shutdown
- statistics: nodes, links, down/up, nomber of node which can reach edge
- store topologies in: JSON / redis / memory
    > can be one mode only, or memory + another method
- authentication
- support east horizontal links
- caching with cache-control
- provide an UI
- UI: provide a diagram of a node
- UI: provide a diagram of a node paths to edges
