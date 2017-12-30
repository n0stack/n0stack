from typing import List, Dict  # NOQA

from n0core.processor import Processor
from n0core.processor import IncompatibleMessage
from n0core.message import Message  # NOQA
from n0core.message import MessageType
from n0core.message.notification import Notification
from n0core.message.notification import Event
from n0core.target import Target  # NOQA
from n0core.gateway import Gateway  # NOQA


class Agent(Processor):
    """Agent is a processor which apply resources with targets.

    1. Receive a message from gateway.
    2. Apply resource with target.
    3. Send a result message to gateway.

    Args:
        target:
        model_types:
        notification:

    Exapmle:
        >>> agent = Agent(notification, [flat_network_target])

    TODO:
        - 引数のnotificationはわかりにくい
    """

    def __init__(self,
                 notification,  # type: Gateway
                 targets=[]     # type: List[Target]
                 ):
        # type: (...) -> None
        self.__notification = notification
        self.__targets = {}  # type: Dict[str, Target]

        for t in targets:
            self.__targets[t.support_model] = t

    @property
    def targets(self):
        # type: () -> Dict[str, Target]
        return self.__targets

    def proccess(self, message):
        # type: (Notification) -> None
        if message.type is not MessageType.NOTIFICATION:
            raise IncompatibleMessage
        if message.model.type not in self.targets.keys():
            raise IncompatibleMessage
        if not message.is_succeeded:
            raise IncompatibleMessage

        model, is_succeeded, description = self.targets[message.model.type].apply(message.model)
        notification = Notification(spec_id=message.spec_id,
                                    model=model,
                                    event=Event.APPLIED,
                                    is_succeeded=is_succeeded,
                                    description=description)

        self.__notification.send(notification)
