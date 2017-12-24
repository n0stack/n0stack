from n0core.processor import Processor
from n0core.processor import IncompatibleMessage
from n0core.message import Message
from n0core.message.notify import Notify


class Agent(Processor):
    def __init__(self,
                 target,       # type: Target
                 model_types,  # type: List[str]
                 notify        # type: Repository
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

        s, d = self.__target.apply(message.model)
        n = Notify(spec_id=message.spec_id,
                   model=message.model,
                   event=Notify.EVENTS.APPLIED,
                   succeeded=s,
                   description=d)

        self.__notify.send(n)
