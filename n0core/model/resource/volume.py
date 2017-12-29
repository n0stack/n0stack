from os.path import join
from enum import Enum
from typing import Dict, List  # NOQA

from n0core.model import Model
from n0core.model import _Dependency # NOQA


class Volume(Model):
    """Volume manage persistent volume resource.

    Example:

        .. code-block:: yaml

        type: resource/volume/file
        id: 486274b2-49e4-4bcd-a60d-4f627ce8c041
        state: allocated
        name: hogehoge
        size: 10 * 1024 * 1024 * 1024
        url: file:///data/hoge

    STATES:
        ALLOCATED: Allocate volume size and share volume.
        DELETED: Delete volume resource, but not delete data in volume.
        DESTROYED: Destroy data in volume.

    Meta:

    Labels:

    Property:

    Args:
        id: UUID
        type:
        state:
        name: Name of volume.
        size: Size of volume.
        url: URL which is sharing like file:///data/hoge and nfs://hobge/data/hoge
        meta:
        dependencies: List of dependency to
    """

    STATES = Enum("STATES", ["ALLOCATED", "DELETED", "DESTROYED"])

    def __init__(self,
                 id,              # type: str
                 type,            # type: str
                 state,           # type: Enum
                 name,            # type: str
                 size,            # type: int
                 url,             # type: str
                 meta={},         # type: Dict[str, str]
                 dependencies=[]  # type: List[_Dependency]
                 ):
        # type: (...) -> None
        super().__init__(id=id,
                         type=join("resource/volume", type),
                         state=state,
                         name=name,
                         meta=meta,
                         dependencies=dependencies)

        self.__size = size
        self.url = url

    @property
    def size(self):
        # type: () -> int
        return self.__size
