from n0core.message import Message
from n0core.gateway import Gateway  # NOQA
from n0core.repository import Repository  # NOQA
from n0core.processor import Processor
from n0core.processor import IncompatibleMessage


class Aggregator(Processor):
    """Aggregator is a processor which store messages.

    1. Receive a message from gateway.
    2. Store messages to repository to provide repository functions.

    Args:
        repository: Data store to store result.

    Exaples:
    """

    def __init__(self, repository):
        # type: (Repository) -> None
        super().__init__()
        self.__repository = repository

    def process(self, message):
        # type: (Message) -> None
        if message.type is not Message.TYPES.NOTIFICATION:
            raise IncompatibleMessage

        self.__repository.store(message)
