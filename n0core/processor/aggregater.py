from n0core.message import Message
from n0core.gateway import Gateway  # NOQA
from n0core.repository import Repository  # NOQA
from n0core.processor import Processor
from n0core.processor import IncompatibleMessage


class Aggregator(Processor):
    def __init__(self, repository):
        # type: (Repository) -> None
        super().__init__()
        self.__repository = repository

    def process(self, message):
        # type: (Message) -> None
        if message.type is not Message.TYPES.NOTIFY:
            raise IncompatibleMessage

        self.__repository.store(message)
