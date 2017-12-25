from typing import Any, List  # NOQA
from enum import Enum  # NOQA

from n0core.model import Model  # NOQA
from n0core.message import Message  # NOQA
from n0core.message.notify import Notify


class Repository:
    def __init__(self):
        # type: () -> None
        pass

    def read(self,
             id,                           # type: str
             *,
             event=Notify.EVENTS.APPLIED,  # type: Enum
             depth=0                       # type: int
             ):
        # type: (...) -> Model
        """
        `read` can get model by id.

        Args:
            id: Model ID such as uuid.
            event: Notify event such as "APPLIED" and "SCHEDULED".
            depth: Depth of model dependency.
                   For example, "VM -> Volume" is 1, "VM" is 0, and "VM -> Volume -> Volume agent" is 2.

        Return:
            Model on event which is setted models until depth.

        Example:
            >>> m = r.read("...", event="APPLIED", depth=1)
            >>> m.dependencies -> not None
            >>> m.dependencies.model.dependencies -> None
        """
        raise NotImplementedError

    def schedule(self, model, ids):
        # type: (Model, List[str]) -> Model
        """
        `schedule` is needed to implement *in the future*.

        Args:
            model: Model of necessary to schedule models.
            ids: List of necessary to create models.

        Return: Model which is attached scheduled agent model.
        """
        pass

    def store(self, message):
        # type: (Message) -> None
        """
        `store` store message to provide query methods like Repository.read and Repository.schedule.

        Args:
            message: Message to store.
                     Model on the top is only stored.
        """
        raise NotImplementedError
