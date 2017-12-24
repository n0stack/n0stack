from enum import Enum
from typing import Any  # NOQA

from n0core.message import Message
from n0core.model import Model  # NOQA


class Notify(Message):
    """
    DO NOT REUSE INSTANCE.
    """

    EVENTS = Enum("EVENTS", ["SCHEDULED", "APPLIED"])

    def __init__(self,
                 spec_id,     # type: str
                 model,       # type: Model
                 event,       # type: Any
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
        # type: () -> Model
        return self.__model

    @property
    def event(self):
        # type: () -> Any
        return self.__event

    @property
    def succeeded(self):
        # type: () -> bool
        return self.__succeeded

    @property
    def description(self):
        # type: () -> str
        return self.__description
