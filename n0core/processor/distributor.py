from n0core.message import Message
from n0core.message.notification import Notification
from n0core.message.spec import Spec  # NOQA
from n0core.processor import Processor
from n0core.processor import IncompatibleMessage
from n0core.gateway import Gateway  # NOQA
from n0core.repository import Repository  # NOQA
from n0core.model import Model  # NOQA


class Distributor(Processor):
    NOTIFICATION_EVENT = Notification.EVENTS.SCHEDULED

    def __init__(self, repository, notification):
        # type: (Repository, Gateway) -> None
        super().__init__()
        self.__repository = repository
        self.__notification = notification

    def applied(self, model):
        # type: (Model) -> bool
        m = self.__repository.read(model.id)

        if m:
            return True
        else:
            return False

    def applied_all(self, model):
        # type: (Model) -> bool
        ms = self.__repository.read(model.id, depth=1)
        ids = map(lambda d: d.model.id, ms.dependencies)

        for i in map(lambda d: d.model.id, model.dependencies):
            if i not in ids:
                return False

        return True

    def process(self, message):
        # type: (Spec) -> None
        if message.type is not Message.TYPES.NOTIFICATION:
            raise IncompatibleMessage

        for m in message.models:
            if self.applied(m):
                continue

            # not scheduled
            if not m.depend_on("resource/hosted"):
                n = Notification(spec_id=message.spec_id,
                                 model=m,
                                 event=self.NOTIFICATION_EVENT,
                                 succeeded=False,
                                 description="not scheduled on your hand.")
                self.__notification.send(n)

            if not self.applied_all(m):
                continue

            a = m.depend_on("resource/hosted")[0].model  # このlabelはfixする必要がある
            n = Notification(spec_id=message.spec_id,
                             model=m,
                             event=self.NOTIFICATION_EVENT,
                             succeeded=True,
                             description="")

            self.__notification.send_to(n, a)
            self.__notification.send(n)
