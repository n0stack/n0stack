from n0core.lib.message import Message


class Repository(object):
    def __init__(self):
        pass

    def read(self,
             id,        # type: str
             *,
             event,     # type: str
             recursive  # type: int
             ):
        # type (...) -> Model
        """
        Example:
            >>> m = r.read("...", event="APPLIED", recursive=1)
            >>> m.dependencies -> not None
            >>> m.dependencies.model.dependencies -> None
        """
        pass

    def store(self, message):
        # type: (Message) -> None
        pass
