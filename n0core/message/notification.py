from enum import Enum
from typing import Any  # NOQA

from n0core.message import Message
from n0core.message import MessageType
from n0core.model import Model  # NOQA


class Event(Enum):
    SCHEDULED = 1
    APPLIED = 2


class Notification(Message):
    """Notification is sent to agent and aggregater to notify a Model.

    Args:
        spec_id: ID to distinguish spec as a user request.
        model: Model that the top of it will be created.
        annotations: Options as scheduling hint and etc.
        event:
        is_succeeded:
        description:

    Example:
        >>> from n0core.model import Model
        >>> m = Model(...)
        >>> Spec(spec_id="ba6f8ced-c8c2-41e9-98d0-5c961dff6c9cf",
                 model=m,
                 event=Event.SCHEDULED,
                 is_succeeded=True,
                 description="Succeeded to schedule {}.".format(m.id))
    """

    def __init__(self,
                 spec_id,     # type: str
                 model,       # type: Model
                 event,       # type: Any
                 is_succeeded,   # type: bool
                 description  # type: str
                 ):
        # type: (...) -> None
        super().__init__(spec_id, MessageType.NOTIFICATION)

        self.__model = model
        self.__event = event
        self.__is_succeeded = is_succeeded
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
    def is_succeeded(self):
        # type: () -> bool
        return self.__is_succeeded

    @property
    def description(self):
        # type: () -> str
        return self.__description
