"""Script to generate topology examples.

Note:
  Written in Python because why not :)
  Love for both Go and Python!
"""
import json
from ipaddress import ip_network
from typing import Generator

PODS = 20
RACKS_PER_POD = 20


def _gen_interco() -> Generator:
    root = ip_network("10.0.0.0/8")
    subnets = root.subnets(23)

    while True:
        yield next(subnets)


gen = _gen_interco()


def create_node(hostname: str, role: str, layer: str) -> dict:
    """Create a node."""
    return {"hostname": hostname, "role": role, "status": True, "layer": layer}


def create_link(south: dict, north: dict) -> dict:
    """Create a link."""
    subnet = next(gen)
    hosts = subnet.hosts()

    return {
        "south_node": south["hostname"],
        "north_node": north["hostname"],
        "status": True,
        "uid": f"{next(hosts)}->{next(hosts)}",
    }


def create_all_nodes() -> list[dict]:
    """Create nodes in the graph."""
    nodes = []
    for i in range(1, PODS + 1):
        # ToR
        for j in range(RACKS_PER_POD):
            hostname = f"tor-{j:02d}-{i:02d}"
            node = create_node(hostname, "tor", 1)
            nodes.append(node)

        # fabric
        for j in range(1, 5):
            hostname = f"fabric-{j}-{i:02d}"
            node = create_node(hostname, "fabric", 2)
            nodes.append(node)

    # spine
    for i in range(1, 5):
        for j in range(1, 5):
            hostname = f"spine-{i}{j}"
            node = create_node(hostname, "spine", 3)
            nodes.append(node)

    # edges
    for i in range(2):
        node = create_node(f"edge-{i}", "edge", 4)
        nodes.append(node)

    return nodes


def create_all_links(nodes: dict) -> list[dict]:
    """Create all links between nodes."""
    WEIGHT = {"tor": 1, "fabric": 2, "spine": 3, "edge": 4}

    links = []

    for south in nodes:
        for north in nodes:
            if WEIGHT[south["role"]] + 1 != WEIGHT[north["role"]]:
                continue

            if south["role"] == "tor":
                pod = south["hostname"][-2:]
                if north["hostname"].endswith(f"-{pod}"):
                    links.append(create_link(south, north))
                    links.append(create_link(south, north))

            if south["role"] == "fabric":
                plane = south["hostname"][-4]
                if north["hostname"][-2] == plane:
                    links.append(create_link(south, north))
                    links.append(create_link(south, north))
                    links.append(create_link(south, north))
                    links.append(create_link(south, north))

            if south["role"] == "spine":
                links.append(create_link(south, north))
                links.append(create_link(south, north))
                links.append(create_link(south, north))
                links.append(create_link(south, north))

    return links


def main() -> None:
    """Generate topology.json file."""
    nodes = create_all_nodes()
    links = create_all_links(nodes)

    with open("topology.json", "w", encoding="utf-8") as fd:
        json.dump({"nodes": nodes, "links": links}, fd, indent=2)


if __name__ == "__main__":
    main()
