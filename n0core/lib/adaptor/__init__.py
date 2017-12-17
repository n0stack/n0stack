from typing import Dict, Any

from n0core.lib.message import Message


class Adapter:
    """
    Adapters provide presentation layer.
    """

    def receive(self):
        # type: () -> Message
        raise NotImplementedError

    def send(self, message):
        # type: (Message) -> None
        """
        This method send message to default destination.

        """
        raise NotImplementedError
