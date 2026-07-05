from src.domain.models import VacancyORM
from src.domain.schemas import UserInfo, VacancyResponseSchema
from src.domain.types import UserRole, VacancyStatus


def has_right_to_vacancy(vacancy: VacancyORM | VacancyResponseSchema, user_info: UserInfo | None) -> bool:
    if vacancy.status == VacancyStatus.PUBLISHED or vacancy.status == VacancyStatus.CLOSED:
        return True

    if user_info is None:
        return False

    return user_info.user_role == UserRole.ADMIN or vacancy.author_id == user_info.user_id


def is_vacancy_id_valid(vacancy_id: int) -> bool:
    return vacancy_id > 0
