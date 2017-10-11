from typing import Any, Callable, Dict, Optional  # NOQA

import n0core.lib.proto
from n0core.lib.messenger import Messenger  # NOQA
from n0core.porter.type import PorterType  # NOQA
from n0core.porter.dhcp.dnsmasq import Dnsmasq  # NOQA


class PorterHandler(object):
    """Messaging queue handler to manage networks on VM hosts.

    Args:
        type: Set PorterType instance what you want to use.
    """

    def __init__(self, type):
        # type: (PorterType) -> None
        self.type = type

    def message_handler(self, messenger):
        # type: (Messenger) -> None
        """Call this when getting message from queue.

        Args:
            messenger: Set Messenger class to send next.
        """
        recv_msg = messenger.inner_msg
        func = getattr(self, recv_msg.__class__.__name__)
        func(recv_msg)

    def __getattr__(self, message_type):
        # type: (str) -> Optional[Callable[[Any], None]]
        """This method call methods named protobuf message.

        Args:
            message_type: Set protobuf sub-message type.

        Returns:
            A method having same protobuf message on PorterType class.

            When the message is undefined on protobuf, return None.
        """
        if hasattr(n0core.lib.proto, message_type):
            return self.__get_type_method(message_type)
        raise AttributeError

    @classmethod
    def __get_type_method(cls, message_type):
        # type: (str) -> Callable[[Any], None]
        f = getattr(cls.type, message_type)  # type: Callable[[Any], None]
        return f
