from n0core.lib.message import Message
from n0core.lib.gateway import Gateway  # NOQA
from n0core.lib.repository import Repository  # NOQA
from n0core.lib.processor import Processor
from n0core.lib.processor import IncompatibleMessage


class Aggregator(Processor):
    def __init__(self, incoming, repository):
        # type: (Gateway, Repository) -> None
        super().__init__(incoming)
        self.__repository = repository

    def process(self, message):
        # type: (Message) -> None
        if message.type is not Message.TYPES.NOTIFY:
            raise IncompatibleMessage

        self.__repository.store(message)
