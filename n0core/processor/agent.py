from typing import List  # NOQA

from n0core.processor import Processor
from n0core.processor import IncompatibleMessage
from n0core.message import Message
from n0core.message.notify import Notify
from n0core.target import Target  # NOQA
from n0core.gateway import Gateway  # NOQA


class Agent(Processor):
    def __init__(self,
                 target,       # type: Target
                 model_types,  # type: List[str]
                 notify        # type: Gateway
                 ):
        # type: (...) -> None
        self.__target = target
        self.__model_types = model_types
        self.__notify = notify

    def proccess(self, message):
        # type: (Notify) -> None
        if message.type is not Message.TYPES.NOTIFY:
            raise IncompatibleMessage
        if message.model.type not in self.__model_types:
            raise IncompatibleMessage
        if not message.succeeded:
            raise IncompatibleMessage

        m, s, d = self.__target.apply(message.model)
        n = Notify(spec_id=message.spec_id,
                   model=m,
                   event=Notify.EVENTS.APPLIED,
                   succeeded=s,
                   description=d)

        self.__notify.send(n)
