from typing import List, Dict  # NOQA

from n0core.processor import Processor
from n0core.processor import IncompatibleMessage
from n0core.message import Message
from n0core.message.notification import Notification
from n0core.target import Target  # NOQA
from n0core.gateway import Gateway  # NOQA


class Agent(Processor):
    def __init__(self,
                 model_types,   # type: List[str]
                 notification,  # type: Gateway
                 target={}      # type: Dict[str, Target]
                 ):
        # type: (...) -> None
        self.__model_types = model_types
        self.__notification = notification
        self.__target = target

    def add_target(self, type, target):
        # type: (str, Target) -> None
        self.__target[type] = target

    def proccess(self, message):
        # type: (Notification) -> None
        if message.type is not Message.TYPES.NOTIFICATION:
            raise IncompatibleMessage
        if message.model.type not in self.__model_types:
            raise IncompatibleMessage
        if not message.is_succeeded:
            raise IncompatibleMessage

        model, is_succeeded, description = self.__target[message.model.type].apply(message.model)
        notification = Notification(spec_id=message.spec_id,
                                    model=model,
                                    event=Notification.EVENTS.APPLIED,
                                    is_succeeded=is_succeeded,
                                    description=description)

        self.__notification.send(notification)
