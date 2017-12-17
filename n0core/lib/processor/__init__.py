from typing import Any, Optional, Dict, List  # NOQA

from n0core.lib.adaptor import Adapter  # NOQA
from n0core.lib.message import Message  # NOQA


class Processor:
    """
    Processor provide logic layer.

    Input and output data is only header(notify and spec) and objects.
    """

    incoming = None  # type: Optional[Adaptor]

    def process(self, message):
        # type: (Message) -> None
        raise NotImplementedError

    def run(self):
        # type: () -> None
        while True:
            nm = self.incoming.receive()
            self.process(nm)
