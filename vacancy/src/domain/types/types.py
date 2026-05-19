from typing import Any, Final, TypeAlias

from pydantic_core import CoreSchema, core_schema

Money: TypeAlias = int
Year: TypeAlias = int


class UnsetValue:
    def __lt__(self, other):
        return True

    def __gt__(self, other):
        return False

    @classmethod
    def __get_pydantic_core_schema__(cls, source_type: Any, handler: Any) -> CoreSchema:
        return core_schema.is_instance_schema(cls)


UNSET_VALUE: Final[UnsetValue] = UnsetValue()
