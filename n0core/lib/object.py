from uuid import uuid4
from typing import Dict, List  # NOQA


class Object(dict):
    """
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
                 dependencies=[]  # List[Dependency]
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


class Dependency():
    """
    Example:
        >>> new_vm = Object("resource/vm/kvm", "running")
        >>> new_disk = Object("resource/volume/local", "claimed")
        >>> new_dependency = Dependency(new_disk, "n0stack/n0core/resource/vm/attachments")
        >>> new_vm.dependencies.append(new_dependency)

    TODO:
        - labelを書き込み可能にするか否か
    """
    def __init__(self,
                 object,      # type: Object
                 label,       # type: str
                 property={}  # type: Dict[str, str]
                 ):
        self.__object = object,
        self.__label = label
        self.property = property

    @property
    def object(self):
        return self.__object

    @property
    def label(self):
        return self.__label