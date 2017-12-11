from typing import Any, Optional, Dict, List  # NOQA

from n0core.lib.processor import Processor


class NotifyProcessor(Processor):
    """
    Example:
        >>> class Aggregator(NotifyProcessor):
        >>>     @on_success
        >>>     def store(self, message):
        >>>         ...
        >>>
        >>>     @on_failure
        >>>     def not_store(self, message):
        >>>         return None
        >>>
        >>> a = Aggregator()
        >>> a.incoming = IncomingPulsar("pulsar://...")
        >>> a.outgoing.append = OutgoingGremlin("ws://...")
        >>> a.run()
    """

    def processing(self, message):
        # type: (Dict[str, Any]) -> Dict[str, Any]
        message = self.pre_proccess(message)

        if message["succeeded"]:
            message = self.on_success(message)
        else:
            message = self.on_failure(message)

        message = self.post_process(message)

        return message

    # these methods will be decorator
    def pre_proccess(self, message):
        # type: (Dict[str, Any]) -> Dict[str, Any]
        return message

    def post_process(self, message):
        # type: (Dict[str, Any]) -> Dict[str, Any]
        return message

    def on_success(self, message):
        # type: (Dict[str, Any]) -> Dict[str, Any]
        return message

    def on_failure(self, message):
        # type: (Dict[str, Any]) -> Dict[str, Any]
        return message
