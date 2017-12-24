from enum import Enum


class Message:
    TYPES = Enum("TYPES", ["NOTIFY", "SPEC"])

    def __init__(self, spec_id, type):
        # type: (str, Enum) -> (None)
        self.__spec_id = spec_id
        self.__type = type

    @property
    def type(self):
        return self.__type

    @property
    def spec_id(self):
        return self.__spec_id
