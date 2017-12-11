from typing import Dict, Any

from n0core.lib.adaptor import Adapter

class IncomingAdapter(Adapter):
    def receive(self):
        # type: () -> Dict[str, Any]
        raise NotImplementedError
