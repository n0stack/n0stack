from typing import Dict, Any  # NOQA

from n0core.lib.message import Message  # NOQA


class Gateway:
    """
    Adapters provide presentation layer.
    """

    def receive(self):
        # type: () -> Message
        raise NotImplementedError

    def send(self, message):
        # type: (Message) -> None
        """
        "send" send message to default destination.

        """
        raise NotImplementedError

    def send_to(self, message, model):
        # type: (Message, Model) -> None
        raise NotImplementedError
