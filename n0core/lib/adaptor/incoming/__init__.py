from typing import Dict, Any

from n0core.lib.adaptor import Adaptor

class IncomingAdaptor(Adaptor):
    def receive(self):
        # type: () -> Dict[str, Any]
        raise NotImplementedError
