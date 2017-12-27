from typing import Dict, List  # NOQA

from n0core.message import Message
from n0core.model import Model  # NOQA


class Spec(Message):
    """Spec is sent from API to scheduler to propagate Models.

    Args:
        spec_id: ID to distinguish spec as a user request.
        models: Models that the top of Model will be created.
        annotations: Options as scheduling hint and etc.

    Example:
        >>> from n0core.model import Model
        >>> m1 = Model(...)
        >>> m2 = Model(...)
        >>> Spec(spec_id="ba6f8ced-c8c2-41e9-98d0-5c961dff6c9cf",
                 models=[m1, m2])
    """

    def __init__(self,
                 spec_id,        # type: str
                 models,         # type: List[Model]
                 annotations={}  # type: Dict[str, str]
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
