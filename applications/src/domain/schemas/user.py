from dataclasses import dataclass
from enum import Enum
from uuid import UUID


class UserRole(Enum):
    UNSPECIFIED = 0
    ADMIN = 1
    EMPLOYEE = 2
    EMPLOYER = 3


@dataclass(slots=True, frozen=True, kw_only=True)
class UserInfo:
    user_role: UserRole
    user_id: UUID
    verificated: bool
