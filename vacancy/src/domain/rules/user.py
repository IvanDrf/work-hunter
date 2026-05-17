from uuid import UUID

from src.domain.schemas.user import UserInfo, UserRole


def is_user_admin(user_info: UserInfo) -> bool:
    return user_info.user_role == UserRole.ADMIN


def is_user_employer(user_info: UserInfo) -> bool:
    return user_info.user_role == UserRole.EMPLOYER


def is_user_vacancy_author(author_id: UUID, user_id: UUID) -> bool:
    return author_id == user_id
