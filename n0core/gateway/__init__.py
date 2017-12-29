from typing import Optional  # NOQA

from n0core.message import Message  # NOQA
from n0core.model import Model  # NOQA


class Gateway:
    """
    Gateway provide methods of incoming or outgoing Messages with other services.
    """

    def receive(self):
        # type: () -> Message
        raise NotImplementedError

    def send(self, message, destination=None):
        # type: (Message, Optional[Model]) -> None
        """
        "send" send message to default destination.

        """
        raise NotImplementedError
