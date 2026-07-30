from pkg.vacancy_api.vacancy_pb2 import VacancyStatus as PKGVacancyStatus

from src.core.exc import ArgumentError
from src.domain.types import VacancyStatus


def vacancy_status_dto(request) -> VacancyStatus:
    statuses = {
        PKGVacancyStatus.MODERATING: VacancyStatus.MODERATING,
        PKGVacancyStatus.PUBLISHED: VacancyStatus.PUBLISHED,
        PKGVacancyStatus.CLOSED: VacancyStatus.CLOSED,
        PKGVacancyStatus.DELETED: VacancyStatus.DELETED,
    }

    if request.status not in statuses:
        raise ArgumentError(f"invalid vacancy status in request, {request.status=}")

    return statuses[request.status]
