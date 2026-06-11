from pkg.vacancy_api.vacancy_pb2 import VacancyStatus as PKGVacancyStatus
from pytest import mark

from src.api.dto.status import vacancy_status_dto
from src.domain.types.enums import VacancyStatus


@mark.parametrize(
    "pkg_status", [PKGVacancyStatus.CLOSED, PKGVacancyStatus.PUBLISHED, PKGVacancyStatus.DELETED, PKGVacancyStatus.MODERATING]
)
def test_vacancy_status_dto(pkg_status: PKGVacancyStatus) -> None:
    status = vacancy_status_dto(pkg_status)

    match pkg_status:
        case PKGVacancyStatus.CLOSED:
            assert status == VacancyStatus.CLOSED

        case PKGVacancyStatus.PUBLISHED:
            assert status == VacancyStatus.PUBLISHED

        case PKGVacancyStatus.MODERATING:
            assert status == VacancyStatus.MODERATING

        case PKGVacancyStatus.DELETED:
            assert status == VacancyStatus.DELETED
