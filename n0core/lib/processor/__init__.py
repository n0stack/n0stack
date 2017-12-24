from typing import Any, Optional, Dict, List  # NOQA

from n0core.lib.adaptor import Adapter  # NOQA
from n0core.lib.message import Message  # NOQA


class IncompatibleMessage(Exception):
    pass


class Processor:
    """
    Processor provide logic layer; in clean architecture, this is similar to UseCase.

    Input and output data is only header(notify and spec) and objects.
    """

    def __init__(self, incoming):
        # type: (Gateway) -> None
        self.__incoming = incoming

    def init(self):
        pass

    def process(self, message):
        # type: (Message) -> None
        raise NotImplementedError

    @incoming.handling
    def handler(self, message):
        try:
            self.process(message)
        except IncompatibleMessage as identifier:
            pass
        except FinishProcess:
            pass
