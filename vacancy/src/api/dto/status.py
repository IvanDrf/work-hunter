from pkg.vacancy_api.vacancy_pb2 import VacancyStatus as PBVacancyStatus

from src.domain.types.enums import VacancyStatus


def vacancy_status_dto(status: PBVacancyStatus) -> VacancyStatus:
    return VacancyStatus(status)
