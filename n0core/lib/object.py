from typing import Dict, List  # NOQA


class Object(dict):
    """

    TODO:
        - dependencyの2重定義ができないようにしたい
    """
    def __init__(self,
                 id,              # str
                 type,            # str
                 state,           # str
                 meta={},         # Dict[str, str]
                 dependencies=[]  # List[Dependency]
                 ):
        # type: (...) -> None
        self.__id = id
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