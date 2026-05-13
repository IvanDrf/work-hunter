from typing import Final, TypeAlias

Money: TypeAlias = int
Year: TypeAlias = int


class UnsetValue:
    def __lt__(self, other):
        return True

    def __gt__(self, other):
        return False


UNSET_VALUE: Final[UnsetValue] = UnsetValue()
