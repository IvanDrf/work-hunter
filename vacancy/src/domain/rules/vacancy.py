from pkg.vacancy_api.vacancy_pb2 import VacancyInfo

from src.core.exc.argument import ArgumentError


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
