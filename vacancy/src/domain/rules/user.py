from src.domain.schemas.user import UserInfo, UserRole


def is_user_admin(user_info: UserInfo) -> bool:
    return user_info.user_role == UserRole.ADMIN


def is_user_employer(user_info: UserInfo) -> bool:
    return user_info.user_role == UserRole.EMPLOYER


def is_user_vacancy_author(vacancy_id: int, user_info: UserInfo) -> bool:
    return vacancy_id == user_info.user_id
