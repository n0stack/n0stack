from n0core.message import Message  # NOQA
from n0core.model import Model  # NOQA


class Gateway:
    """Gateway provide methods of incoming or outgoing Messages with other services.

    TODO:
        仕様が決まっていない
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

    def send_to(self, message, destination):
        # type: (Message, Model) -> None
        """
        Args:
            message:
            destination: destination
        """
        raise NotImplementedError
