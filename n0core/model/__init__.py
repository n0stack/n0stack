from typing import Dict, List  # NOQA
from enum import Enum  # NOQA


class Model:
    """Model is mapped to express graph data structure with _Dependency.

    About meta:
        There are a lot of way to collaborate other service.

        #### Example 1: n0gateway

        以下のようなリクエストをユーザーが行うと `9c2476ab-dc1e-4904-b8a4-6d991fdc7770` のUUIDに関連付けられているLBサービスにportが参加する。

        ```yaml
        type: resuorce/nic
        state: running
        meta:
        n0stack/n0gateway/join: 9c2476ab-dc1e-4904-b8a4-6d991fdc7770
        ```

        n0gatewayとしては `/api/spec` を監視していればサービスディスカバリを用意に実装することができる。

    Args:
        id: UUID  default: generate uuid
        type:
        state: 
        meta: 
        dependencies: List of dependency to 

    Example:
        >>> new_disk = Model("resource/volume/local", "claimed")
        >>> new_disk["size"] = 100 * 1024 * 1024 * 1024
        >>> new_disk.meta["n0stack/n0core/resource/vm/boot_priority"] = "1"

    TODO:
        - dependencyの2重定義ができないようにしたい
    """

    def __init__(self,
                 id,              # type: str
                 type,            # type: str
                 state,           # type: Enum
                 name,            # type: str
                 meta={},         # type: Dict[str, str]
                 dependencies=[]  # type: List[_Dependency]
                 ):
        # type: (...) -> None
        self.__id = id
        self.__type = type
        self.state = state
        self.__name = name
        self.meta = meta
        self.dependencies = dependencies

    @property
    def id(self):
        # type: () -> str
        return self.__id

    @property
    def type(self):
        # type: () -> str
        return self.__type

    @property
    def name(self):
        # type: () -> str
        return self.__name

    def depend_on(self, label):
        # type: (str) -> List[_Dependency]
        return [d for d in self.dependencies if d.label == label]

    def add_dependency(self,
                       model,       # type: Model
                       label,       # type: str
                       property={}  # type: Dict[str, str]
                       ):
        # type: (...) -> None
        d = _Dependency(model, label, property)

        for i, v in enumerate(self.dependencies):
            if v.model.id == d.model.id:
                self.dependencies.pop(i)

        self.dependencies.append(d)


class _Dependency:
    """
    Example:
        >>> new_vm = Model("resource/vm/kvm", "running")
        >>> new_disk = Model("resource/volume/local", "claimed")
        >>> new_vm.add_dependency(new_disk, "resource/vm/attachments")

    TODO:
        - labelを書き込み可能にするか否か
    """
    def __init__(self,
                 model,       # type: Model
                 label,       # type: str
                 property={}  # type: Dict[str, str]
                 ):
        # type: (...) -> None
        self.__model = model
        self.__label = label
        self.property = property

    @property
    def model(self):
        # type: () -> Model
        return self.__model

    @property
    def label(self):
        # type: () -> str
        return self.__label
