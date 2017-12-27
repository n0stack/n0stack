from typing import List  # NOQA

from n0core.processor import Processor
from n0core.processor import IncompatibleMessage
from n0core.message import Message
from n0core.message.notification import Notification
from n0core.target import Target  # NOQA
from n0core.gateway import Gateway  # NOQA


class Agent(Processor):
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
        if not message.succeeded:
            raise IncompatibleMessage

        model, succeeded, description = self.__target.apply(message.model)
        notification = Notification(spec_id=message.spec_id,
                                    model=model,
                                    event=Notification.EVENTS.APPLIED,
                                    succeeded=succeeded,
                                    description=description)

        self.__notification.send(notification)
