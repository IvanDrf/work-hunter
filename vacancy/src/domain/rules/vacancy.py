from typing import Final

from pkg.common.common_pb2 import UserInfo, UserRole
from pkg.vacancy_api.vacancy_pb2 import VacancyInfo

from src.core.exc import ArgumentError
from src.domain.models.vacancy import VacancyORM, VacancyStatus


MIN_VACANCY_ID: Final[int] = 0


def check_vacancy_fields(vacancy: VacancyInfo) -> None:
    if vacancy.title == '':
        raise ArgumentError(
            f'title must be non empty, but given: {vacancy.title}'
        )

    if not _is_salary_valid(vacancy):
        raise ArgumentError(
            f'salary_min must be less or equal to salary_max, but given: {vacancy.salary_min} and {vacancy.salary_max}'
        )

    if not _is_experience_valid(vacancy):
        raise ArgumentError(
            f'experience_min must be less or equal to experience_max, but given: {vacancy.experience_min} and {vacancy.experience_max}'
        )


def has_right_to_vacancy(vacancy: VacancyORM, user_info: UserInfo) -> bool:
    return vacancy.status == VacancyStatus.MODERATING and user_info.role != UserRole.ADMIN and vacancy.author_id != user_info.user_id


def is_vacancy_id_valid(vacancy_id: int) -> bool:
    return vacancy_id < MIN_VACANCY_ID


def _is_salary_valid(vacancy: VacancyInfo) -> bool:
    # if have only min salary
    if vacancy.salary_max == 0:
        return vacancy.salary_min >= 0

    return vacancy.salary_max >= vacancy.salary_min and vacancy.salary_min >= 0


def _is_experience_valid(vacancy: VacancyInfo) -> bool:
    # if have only min experience
    if vacancy.experience_max == 0:
        return vacancy.experience_min >= 0

    return vacancy.experience_max >= vacancy.experience_min and vacancy.experience_min >= 0
