from typing import Dict, List  # NOQA

from n0core.message import Message
from n0core.model import Model  # NOQA


class Spec(Message):
    def __init__(self,
                 spec_id,     # type: str
                 models,      # type: List[Model]
                 annotations  # type: Dict[str, str]
                 ):
        # type: (...) -> None
        super().__init__(spec_id, Message.TYPES.SPEC)

        self.__models = models
        self.__annotations = annotations

    @property
    def models(self):
        # type: () -> List[Model]
        return self.__models

    @property
    def annotations(self):
        # type: () -> Dict[str, str]
        return self.__annotations
