from typing import Any, Final, TypeAlias

from pydantic_core import CoreSchema, core_schema

Money: TypeAlias = int
Year: TypeAlias = int


class SingleTon(type):
    _instances = {}

    def __call__(cls, *args: Any, **kwargs: Any) -> Any:
        if cls not in cls._instances:
            cls._instances[cls] = super(SingleTon, cls).__call__(*args, **kwargs)

        return cls._instances[cls]


class UnsetValue(metaclass=SingleTon):
    def __lt__(self, other):
        return True

    def __gt__(self, other):
        return False

    def __contains__(self, item):
        return False

    @classmethod
    def __get_pydantic_core_schema__(cls, source_type: Any, handler: Any) -> CoreSchema:
        return core_schema.is_instance_schema(cls)


UNSET_VALUE: Final[UnsetValue] = UnsetValue()
