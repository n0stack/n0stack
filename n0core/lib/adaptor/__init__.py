from typing import Dict, Any

class Adapter:
    """
    Adapters provide presentation layer.
    """

    def receive(self):
        # type: () -> Dict[str, Any]
        raise NotImplementedError

    def send(self, message):
        # type: (Dict[str, Any], List[Dict[str, Any]]) -> None
        """
        This method send message to default destination.

        Args:
            header:
            objects:
        """
        raise NotImplementedError
