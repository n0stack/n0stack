from uuid import uuid4

from n0core.model.resource.network import Network
from n0core.model.resource.nic import NIC
from n0core.target.network.flat import Flat


INTERFACE_NAME = ""
NETWORK_ID = "00f967af-0cc6-443e-b2a9-277b158f42f0"

up_network = Network(NETWORK_ID, "flat", Network.STATES.UP, "test_network")
up_network.add_subnet("192.168.0.0/24", "192.168.0.3-192.168.0.127", ["192.168.0.254"], "192.168.0.254")

delete_network = Network(NETWORK_ID, "flat", Network.STATES.DELETED, "test_network", "nflat00f967af")
delete_network.add_subnet("192.168.0.0/24", "192.168.0.3-192.168.0.127", ["192.168.0.254"], "192.168.0.254")

nic = NIC("755cedc3-1fe2-44a8-bc45-93d8edf2ad82", "dhcp", NIC.STATES.ATTACHED, "test_nic", "1234567890abcd", ["192.168.0.1"])
nic.add_dependency(up_network, "n0stack/n0core/resource/port/network")

if __name__ == "__main__":
    f = Flat(INTERFACE_NAME)
    f.apply(up_network)
    # f.apply(nic)
