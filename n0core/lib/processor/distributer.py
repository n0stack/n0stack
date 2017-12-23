from n0core.lib.message import Message
from n0core.lib.message.notify import Notify
from n0core.lib.message.notify import Spec
from n0core.lib.processor import Processor
from n0core.lib.processor import IncompatibleMessage


class Distributer(Processor):
    NOTIFY_EVENT = "SCHEDULED"

    def __init__(self, incoming, repository, notify):
        # type: (Gateway, Gateway) -> None
        super().__init__(incoming)
        self.__repository = repository
        self.__notify = notify

    def applied(self, model):
        # type: (Model) -> boolean
        m = self.__repository.read(model.id, event="APPLIED", recursive=0)

        if m:
            return True
        else:
            return False

    def process(self, message):
        # type: (Spec) -> None
        if message.type is not Message.TYPES.NOTIFY:
            raise IncompatibleMessage

        for m in message.models:
            if self.applied(m):
                continue

            # not scheduled
            if not m.depend_on("n0core/models/hosted"):
                n = Notify(spec_id=message.spec_id,
                           model=m,
                           event=self.NOTIFY_EVENT,
                           succeeded=False,
                           description="")
                self.__notify.send(n)

            for r in filter(lambda r: r.model.type.split("/")[0] == "resource", m.relations):
                if not self.applied(r):
                    continue

            a = m.depend_on("n0core/models/hosted")[0].model  # このlabelはfixする必要がある
            n = Notify(spec_id=message.spec_id,
                       model=m,
                       event=self.NOTIFY_EVENT,
                       succeeded=True,
                       description="")

            try:
                self.__notify.send_to(n, a)
            except:
                n = Notify(spec_id=message.spec_id,
                           model=m,
                           event=self.NOTIFY_EVENT,
                           succeeded=False,
                           description="Fail sending message to agent")

            self.__notify.send(n)
