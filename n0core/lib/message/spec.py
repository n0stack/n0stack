from typing import Dict  # NOQA

from n0core.lib.message import Message
from n0core.lib.object import Object  # NOQA


class Spec(Message):
    def __init__(self,
                 spec_id,     # type: str
                 models,     # type: List[Model]
                 annotations  # type: Dict[str, str]
                 ):
        super().__init__(spec_id, Message.TYPES.SPEC)

        self.__models = models
        self.__annotations = annotations

    @property
    def models(self):
        return self.__models

    @property
    def annotations(self):
        return self.__annotations
