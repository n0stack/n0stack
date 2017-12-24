from enum import Enum

from n0core.message import Message
from n0core.object import Object # NOQA


class Notify(Message):
    """
    DO NOT REUSE INSTANCE.
    """

    EVENTS = Enum("EVENTS", ["SCHEDULED", "APPLIED"])

    def __init__(self,
                 spec_id,     # type: str
                 model,      # type: Model
                 event,       # type: str
                 succeeded,   # type: bool
                 description  # type: str
                 ):
        # type: (...) -> None
        super().__init__(spec_id, Message.TYPES.NOTIFY)

        self.__model = model
        self.__event = event
        self.__succeeded = succeeded
        self.__description = description

    @property
    def model(self):
        return self.__model

    @property
    def event(self):
        return self.__event

    @property
    def succeeded(self):
        return self.__succeeded

    @property
    def description(self):
        return self.__description
