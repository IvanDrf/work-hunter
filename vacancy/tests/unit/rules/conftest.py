from typing import TypeAlias

from pytest import fixture

from src.domain.types.types import UNSET_VALUE, Money, UnsetValue

Salaries: TypeAlias = list[tuple[Money | UnsetValue | None, Money | UnsetValue | None]]


@fixture(scope="module")
def valid_salaries() -> Salaries:
    return [
        (5, 67),
        (0, 0),
        (20, 20),
        (None, None),
        (20, None),
        (None, 30),
        (UNSET_VALUE, 30),
        (30, UNSET_VALUE),
        (UNSET_VALUE, UNSET_VALUE),
    ]


@fixture(scope="module")
def invalid_salaries() -> Salaries:
    return [
        (20, 10),
        (-5, 10),
        (-5, -5),
        (None, -5),
        (-5, None),
        (UNSET_VALUE, -5),
        (-5, UNSET_VALUE),
    ]
