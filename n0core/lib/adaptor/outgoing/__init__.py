from typing import Dict, Any

from n0core.lib.adaptor import Adapter

class OutgoingAdapter(Adapter):
    def send(self, message):
        # type: (Dict[str, Any]) -> None
        raise NotImplementedError
