from typing import Any, Optional, Dict, List  # NOQA

from n0core.message import Message  # NOQA


class IncompatibleMessage(Exception):
    pass


class Processor:
    """
    TODO:
        - どうやってincomingのGatewayからいい感じにデータを取得するか
    """

    def __init__(self):
        # type: () -> None
        pass

    def process(self, message):
        # type: (Message) -> None
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
