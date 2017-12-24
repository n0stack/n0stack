from typing import Tuple  # NOQA

from n0core.target import Target
from n0core.model import Model  # NOQA
from n0core.target.network.dhcp.dnsmasq import Dnsmasq


class Network(Target):
    def __init__(self, bridge_type):
        # type: () -> None
        self.__type = bridge_type  # instance
        self.__dhcp = Dnsmasq.get_created()

    def apply(self, model):
        # type: (Model) -> Tuple[bool, str]
        if model.type.split("/")[0] != "resource":
            return False, "This model is not supported."

        p = self.decode_parameters(model)

        resource_type = model.type.split("/")[1]
        if resource_type == "network":
            if model.state == "up":
                model["bridge"] = b.apply_bridge(model.id, parameter=p)  # vlan idなどをどうやってわたすか

                for s in model["subnets"]:
                    d = self.__dhcp[model.id][s["cidr"]]

                    if d is None:  # dhcp is created
                        self.self.__dhcp[model.id][s["cidr"]] = Dnsmasq(model.id, model["bridge"]))  # Check already exists dnsmasq instance  cidrが一意なはず
                        d.create_dhcp_server(s["cidr"], model["bridge"], s["dhcp"]["range"])

                    else:
                        d.respawn_process(s["dhcp"]["range"])

            elif model.state == "down":
                self.__type.apply_bridge(model["bridge"], state=down, parameter=p)
                for s in model["subnets"]:
                    self.__dhcp[model.id][s["cidr"]].stop_process()

            elif model.state == "deleted":
                self.__type.delete_bridge(model["bridge"])
                for s in model["subnets"]:
                    self.__dhcp[model.id][s["cidr"]].delete_dhcp_server()

        elif resource_type == "port":
            nid = model.depend_on("n0core/port/network")[0].model.id
            d = self.__dhcp[nid][s["cidr"]]

            if model.state == "attached":
                d.add_allowed_host(model["hw_addr"])

                for i in model["ip_addrs"]:
                    d.add_host_entry(model["hw_addr"], i)

            elif model.state == "detached":
                d.delete_host_entry(model["hw_addr"])
                d.delete_allowed_host(model["hw_addr"])

    def decode_parameters(self, model):
        # type: (Model) -> Dict[str, str]
        d = {}

        for k, v in filter(lambda k, v: self.__type.meta_prefix in k, model.meta.items()):
            d[k.split("/")[-1]] = v

        return d
