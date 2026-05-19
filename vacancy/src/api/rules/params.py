from typing import Final, Sequence

from src.core.exc import ArgumentError

MIN_OFFSET: Final[int] = 0

MIN_LIMIT: Final[int] = 5
MAX_LIMIT: Final[int] = 30

MIN_TAGS_AMOUNT: Final[int] = 1
MAX_TAGS_AMOUNT: Final[int] = 2


def validate_limit_offset(limit: int, offset: int) -> None:
    if not is_limit_valid(limit):
        raise ArgumentError(f"limit must be in range ({MIN_LIMIT}, {MAX_LIMIT}), but {limit=}")

    if not is_offset_valid(offset):
        raise ArgumentError(f"offset must be greater than {MIN_OFFSET}, but {offset=}")


def is_offset_valid(offset: int) -> bool:
    return offset >= MIN_OFFSET


def is_limit_valid(limit: int) -> bool:
    return MIN_LIMIT <= limit <= MAX_LIMIT


def is_tags_amount_valid(tags: Sequence) -> bool:
    return MIN_TAGS_AMOUNT <= len(tags) <= MAX_TAGS_AMOUNT
