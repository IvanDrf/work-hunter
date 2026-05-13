from enum import Enum


class RemoteType(Enum):
    OFFICE = 0
    REMOTE = 1
    HYBRID = 2
    ANY = 3


class TimeType(Enum):
    FULL = 0
    PART = 1


class Currency(Enum):
    RUB = 0
    USD = 1
    EUR = 2


class VacancyStatus(Enum):
    MODERATING = 0
    PUBLISHED = 1
    UPDATED = 2
    CLOSED = 3
    DELETED = 4
