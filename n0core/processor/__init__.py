from typing import Any, Optional, Dict, List  # NOQA

from n0core.message import Message  # NOQA


class Processor:
    """
    TODO:
        - どうやってincomingのGatewayからいい感じにデータを取得するか
    """

    def __init__(self):
        # type: () -> None
        pass

    def process(self, message):
        # type: (Any) -> None
        """
        `process` is implimented by each use case.

        Args:
            message: Message which is got by Gateway.
        """
        raise NotImplementedError

    def handler(self, message):
        # type: (Message) -> None
        """
        `handler` wrap `Processor.process` to manage common processes over all like exceptions.

        Gateway will call this.

        Args:
            message: Message which is got by Gateway.
        """
        try:
            self.process(message)
        except IncompatibleMessage as identifier:
            pass


class IncompatibleMessage(Exception):
    """Raise when received not supported message.

    You sent wrong message or forgot the message implementation.

    Args:
        message_type: Set message type defined on protobuf.
    """

    def __init__(self, message_type):
        # type: (str) -> None
        self.message_type = message_type

    def __str__(self):
        # type: () -> str
        return "on {}".format(self.message_type)
