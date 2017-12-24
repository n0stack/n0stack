from typing import Any, List

from n0core.model import Model  # NOQA
from n0core.message import Message  # NOQA
from n0core.message.notify import Notify


class Repository:
    def __init__(self):
        # type: () -> None
        pass

    def read(self,
             id,                          # type: str
             *,
             event=Notify.TYPES.APPLIED,  # type: Any
             depth=0                      # type: int
             ):
        # type (...) -> Model
        """
        Example:
            >>> m = r.read("...", event="APPLIED", depth=1)
            >>> m.dependencies -> not None
            >>> m.dependencies.model.dependencies -> None
        """
        raise NotImplementedError

    def schedule(self, model, ids):
        # type: (Model, List[str]) -> Model
        """
        Args:
            model: Model of necessary to schedule models.
            ids: List of necessary to create models.

        Return: Model which is attached scheduled agent model.
        """
        pass

    def store(self, message):
        # type: (Message) -> None
        raise NotImplementedError
