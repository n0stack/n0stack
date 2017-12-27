from enum import Enum
from typing import Any  # NOQA

from n0core.message import Message
from n0core.model import Model  # NOQA


class Notification(Message):
    """Spec is sent from API to scheduler to propagate Models.

    Args:
        spec_id: ID to distinguish spec as a user request.
        model: Model that the top of it will be created.
        annotations: Options as scheduling hint and etc.
        event:
        succeeded:
        description:

    Example:
        >>> from n0core.model import Model
        >>> m = Model(...)
        >>> Spec(spec_id="ba6f8ced-c8c2-41e9-98d0-5c961dff6c9cf",
                 model=m,
                 event=Notification.EVENTS.SCHEDULED,
                 succeeded=True,
                 description="Succeeded to schedule {}.".format(m.id))
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
        super().__init__(spec_id, Message.TYPES.NOTIFICATION)

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
