from typing import Any, Optional, Dict, List  # NOQA

from n0core.lib.adaptor.incoming import IncomingAdapter
from n0core.lib.adaptor.outgoing import OutgoingAdapter

class Processor:
    incoming = None  # type: Optional[IncomingAdapter]
    outgoing = []  # type: List[OutgoingAdapter]

    def __init__(self):
        # type: () -> None
        pass

    def processing(self, message):
        # type: (Dict[str, Any]) -> Dict[str, Any]
        return message

    def run(self):
        # type: () -> None
        while True:
            nm = self.incoming.receive()
            pm = self.processing(nm)

            for o in self.outgoing:
                o.send(pm)
