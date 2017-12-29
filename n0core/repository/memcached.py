from memcache import Client

from typing import Any, List, cast  # NOQA
from enum import Enum  # NOQA

from n0core.model import Model  # NOQA
from n0core.message import Message  # NOQA
from n0core.message.notification import Notification


class Memcached:
    """
    Args:
        memcached: List of host and path to connect memcached like ["127.0.0.1:11211"].
    """

    def __init__(self, memcached):
        # type: (List[str]) -> None
        super().__init__()

        self.__client = Client(memcached, cache_cas=True)

    def read(self,
             id,                                 # type: str
             *,
             event=Notification.EVENTS.APPLIED,  # type: Enum
             depth=0                             # type: int
             ):
        # type: (...) -> Model
        return self._get_model(id, event.name, depth)

    def _get_model(self, id, event, depth):
        # type: (str, str, int) -> Model
        key = self.get_key(id, event)
        model = self.__client.get(key)  # type: Model

        if depth == 0:
            model.dependencies = []
            return model

        for i, d in enumerate(model.dependencies):
            model.add_dependency(self._get_model(d.model.id, event, depth-1), d.label, d.property)

        return model

    def store(self, message):
        # type: (Message) -> None
        message = cast(Notification, message)
        key = self.get_key(message.event.name, message.model.id)
        self.__client.set(key, message.model)

    def get_key(self, id, event):
        # type: (str, str) -> str
        return "{}-{}".format(event, id)
