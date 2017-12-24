from n0core.lib.message import Message


class Repository:
    def __init__(self):
        pass

    def read(self,
             id,               # type: str
             *,
             event="APPLIED",  # type: str
             depth=0           # type: int
             ):
        # type (...) -> Model
        """
        Example:
            >>> m = r.read("...", event="APPLIED", depth=1)
            >>> m.dependencies -> not None
            >>> m.dependencies.model.dependencies -> None
        """
        raise NotImplementedError

    def store(self, message):
        # type: (Message) -> None
        raise NotImplementedError
