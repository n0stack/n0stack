from uuid import uuid4
from typing import Dict, List  # NOQA


class Model(dict):
    """
    [WARNING] 仕様が変わる可能性があるので、それを考慮して開発をするように!!

    Model and Dependency is mapped to express graph data structure.
    See details in /doc/architecture.

    Example:
        >>> new_disk = Object("resource/volume/local", "claimed")
        >>> new_disk["size"] = 100 * 1024 * 1024 * 1024
        >>> new_disk.meta["n0stack/n0core/resource/vm/boot_priority"] = "1"

    TODO:
        - dependencyの2重定義ができないようにしたい
    """
    def __init__(self,
                 type,            # str
                 state,           # str
                 id="",           # str
                 meta={},         # Dict[str, str]
                 dependencies=[]  # List[_Dependency]
                 ):
        # type: (...) -> None
        if id:
            self.__id = id
        else:
            self.__id = uuid4()

        self.__type = type
        self.state = state
        self.meta = meta
        self.dependencies = dependencies

    @property
    def id(self):
        return self.__id

    @property
    def type(self):
        return self.__type

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
        >>> new_vm.add_dependency(new_disk, "n0stack/n0core/resource/vm/attachments")

    TODO:
        - labelを書き込み可能にするか否か
    """
    def __init__(self,
                 model,       # type: Model
                 label,       # type: str
                 property={}  # type: Dict[str, str]
                 ):
        # type: (...) -> None
        self.__model = model,
        self.__label = label
        self.property = property

    @property
    def model(self):
        return self.__model

    @property
    def label(self):
        return self.__label
