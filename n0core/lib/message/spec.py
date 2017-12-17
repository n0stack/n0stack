from typing import Dict  # NOQA

from n0core.lib.message import Message
from n0core.lib.object import Object  # NOQA


class Spec(Message):
    def __init__(self,
                 spec_id,     # type: stiring
                 objects,     # type: List[Objects]
                 annotations  # type: Dict[string, string]
                 ):
        super().__init__(spec_id, Message.TYPES.SPEC)

        self.__spec = objects
        self.__annotations = annotations

    @property
    def objects(self):
        return self.__spec

    @property
    def annotations(self):
        return self.__annotations
