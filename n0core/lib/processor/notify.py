from typing import Any, Optional, Dict, List  # NOQA

from n0core.lib.processor import Processor


class NotifyProcessor(Processor):
    """
    Example:
        >>> class Aggregator(NotifyProcessor):
        >>>     def on_success(self, message):
        >>>         ...
        >>>
        >>>     def on_failure(self, message):
        >>>         pass
        >>>
        >>> a = Aggregator()
        >>> a.incoming = IncomingPulsar("pulsar://...")
        >>> a.outgoing.append(OutgoingGremlin("ws://..."))
        >>> a.run()
    """

    def process(self, message):
        # type: (Message) -> None
        if message.type != message.NOTIFY:
            raise Exception("Not supported message.")

        message = self.pre_proccess(message)

        if message["succeeded"]:
            message = self.on_success(message)
        else:
            message = self.on_failure(message)

        message = self.post_process(message)

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
