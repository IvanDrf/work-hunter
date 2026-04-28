from typing import Final, Sequence


MIN_OFFSET: Final[int] = 0

MIN_LIMIT: Final[int] = 5
MAX_LIMIT: Final[int] = 30

MIN_TAGS_AMOUNT: Final[int] = 1
MAX_TAGS_AMOUNT: Final[int] = 2


def is_offset_valid(offset: int) -> bool:
    return offset >= MIN_OFFSET


def is_limit_valid(limit: int) -> bool:
    return MIN_LIMIT <= limit <= MAX_LIMIT


def is_tags_amount_valid(tags: Sequence) -> bool:
    return MIN_TAGS_AMOUNT <= len(tags) <= MAX_TAGS_AMOUNT
