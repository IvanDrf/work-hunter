from dataclasses import dataclass
from uuid import UUID

from src.domain.types.enums import UserRole


@dataclass(slots=True, frozen=True, kw_only=True)
class UserInfo:
    user_role: UserRole
    user_id: UUID
    verificated: bool
