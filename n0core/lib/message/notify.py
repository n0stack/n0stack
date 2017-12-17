from enum import Enum

from n0core.lib.message import Message
from n0core.lib.object import Object # NOQA


class Notify(Message):
    """
    DO NOT REUSE.
    """

    EVENTS = Enum("EVENTS", ["SCHEDULED", "APPLIED"])

    def __init__(self,
                 spec_id,     # type: string
                 object,      # type: Object
                 event,       # type: string
                 succeeded,   # type: bool
                 description  # type: string
                 ):
        # type: (...) -> None
        super().__init__(spec_id, Message.TYPES.NOTIFY)

        self.__object = object
        self.__event = self.EVENTS[event]
        self.__succeeded = succeeded
        self.__description = description

    @property
    def object(self):
        return self.__object

    @property
    def event(self):
        return self.__event

    @property
    def succeeded(self):
        return self.__succeeded

    @property
    def description(self):
        return self.__description