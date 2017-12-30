from typing import List  # NOQA

from n0core.processor import Processor
from n0core.processor import IncompatibleMessage
from n0core.message import Message
from n0core.message.notification import Notification
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
 
    TODO:
        - 引数のnotificationはわかりにくい
    """

    def __init__(self,
                 target,       # type: Target
                 model_types,  # type: List[str]
                 notification  # type: Gateway
                 ):
        # type: (...) -> None
        self.__target = target
        self.__model_types = model_types
        self.__notification = notification

    def proccess(self, message):
        # type: (Notification) -> None
        if message.type is not Message.TYPES.NOTIFICATION:
            raise IncompatibleMessage
        if message.model.type not in self.__model_types:
            raise IncompatibleMessage
        if not message.is_succeeded:
            raise IncompatibleMessage

        model, is_succeeded, description = self.__target.apply(message.model)
        notification = Notification(spec_id=message.spec_id,
                                    model=model,
                                    event=Notification.EVENTS.APPLIED,
                                    is_succeeded=is_succeeded,
                                    description=description)

        self.__notification.send(notification)
