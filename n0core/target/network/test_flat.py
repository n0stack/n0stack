from uuid import uuid4

from n0core.model import Model
from n0core.target.network.flat import Flat


INTERFACE_NAME = ""
NETWORK_ID = "00f967af-0cc6-443e-b2a9-277b158f42f0"

up_network = Model("resource/network/flat", "up", id=NETWORK_ID)
up_network["subnets"] = [{
    "cidr": "192.168.0.0/24",
    "dhcp": {
        "range": "192.168.0.3-192.168.0.127",
        "nameservers": [
            "192.168.0.254"
        ],
        "gateway": "192.168.0.254"
    }
}]
down_network = Model("resource/network/flat", "down", id=NETWORK_ID)
down_network["subnets"] = up_network["subnets"]
down_network["bridge"] = "nflat00f967af"
delete_network = Model("resource/network/flat", "deleted", id=NETWORK_ID)
delete_network["subnets"] = up_network["subnets"]
delete_network["bridge"] = "nflat00f967af"

port = Model("resource/port",
             "up")
port["hw_addr"] = "1234567890ab"
port["ip_addrs"] = ["192.168.0.3"]
port.add_dependency(up_network, "n0core/resource/port/network")

if __name__ == "__main__":
    f = Flat(INTERFACE_NAME)
    f.apply(up_network)
    # f.apply(port)
