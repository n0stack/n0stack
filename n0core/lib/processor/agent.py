from n0core.lib.processor import Processor
from n0core.lib.processor import IncompatibleMessage
from n0core.lib.processor import FinishProcess
from n0core.lib.message import Message
from n0core.lib.message.notify import Notify


class Agent(Processor):
    def __init__(self, target, model_types, notify):
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
            raise FinishProcess

        s, d = self.__target.apply(message.model)

        n = Notify(spec_id=message.spec_id,
                   model=message.model,
                   event="APPLIED",
                   succeeded=s,
                   description=d)

        self.__notify.send(n)
