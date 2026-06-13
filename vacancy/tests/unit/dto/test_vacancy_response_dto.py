from datetime import datetime, timezone

from hypothesis import given
from hypothesis import strategies as st
from pkg.vacancy_api.vacancy_pb2 import Currency as PKGCurrency
from pkg.vacancy_api.vacancy_pb2 import RemoteType as PKGRemoteType
from pkg.vacancy_api.vacancy_pb2 import TimeType as PKGTimeType
from pkg.vacancy_api.vacancy_pb2 import VacancyInfo as PKGVacancyInfo

from src.api.dto.vacancy import vacancy_response_dto
from src.domain.schemas import VacancyResponseSchema
from src.domain.types.enums import Currency, RemoteType, TimeType, VacancyStatus
from tests.unit.dto.common import FullVacancyInfo, full_vacancy_info


@given(responses=st.lists(full_vacancy_info(), min_size=5, max_size=5))
def test_vacancy_response_dto(responses: list[FullVacancyInfo]) -> None:
    for resp in responses:
        vacancy = _create_vacancy_response(resp)
        info = vacancy_response_dto(vacancy)

        assert_vacancy_response(vacancy, info)


def assert_vacancy_response(vacancy: VacancyResponseSchema, info: PKGVacancyInfo) -> None:
    assert info.vacancy_id == vacancy.vacancy_id
    assert info.title == vacancy.title
    assert info.description == vacancy.description
    assert info.requirements == vacancy.requirements
    assert info.conditions == vacancy.conditions
    assert info.salary_min == vacancy.salary_min
    assert info.salary_max == vacancy.salary_max
    assert info.currency == vacancy.currency.name
    assert info.experience_min == vacancy.experience_min
    assert info.experience_max == vacancy.experience_max
    assert info.created_at == vacancy.created_at
    assert info.status == vacancy.status.name
    assert info.remote_type == vacancy.remote_type.name
    assert info.time_type == vacancy.time_type.name
    assert info.city == vacancy.city
    assert info.metro == vacancy.metro
    assert info.views == vacancy.views
    assert info.applications_count == vacancy.applications_count
    assert info.tags == vacancy.tags
    assert info.author_name == vacancy.author_name
    assert info.moderated_time == vacancy.moderated_at
    assert info.moderator_comments == vacancy.moderator_comments
    assert info.updated_at == vacancy.updated_at
    assert info.published_at == vacancy.published_at
    assert info.closed_at == vacancy.closed_at


def _create_vacancy_response(resp: FullVacancyInfo) -> VacancyResponseSchema:
    time_types = {
        PKGTimeType.FULL: TimeType.FULL,
        PKGTimeType.PART: TimeType.PART,
    }

    remotes = {
        PKGRemoteType.ANY: RemoteType.ANY,
        PKGRemoteType.HYBRID: RemoteType.HYBRID,
        PKGRemoteType.OFFICE: RemoteType.OFFICE,
        PKGRemoteType.REMOTE: RemoteType.REMOTE,
    }

    currs = {
        PKGCurrency.EUR: Currency.EUR,
        PKGCurrency.USD: Currency.USD,
        PKGCurrency.RUB: Currency.RUB,
    }

    return VacancyResponseSchema(
        title=resp.main.title,
        description=resp.main.description,
        requirements=resp.main.requirements,
        conditions=resp.main.conditions,
        salary_min=resp.salary.salary_min,
        salary_max=resp.salary.salary_max,
        currency=currs[resp.salary.currency],
        city=resp.additional.city,
        metro=resp.additional.metro,
        remote_type=remotes[resp.additional.remote_type],
        time_type=time_types[resp.additional.time_type],
        experience_min=resp.exp.experience_min,
        experience_max=resp.exp.experience_max,
        tags=resp.main.tags,
        vacancy_id=resp.vacancy_id,
        author_name=resp.author.author_name,
        author_id=resp.author.author_id,
        status=resp.stats.status,
        created_at=resp.time.created_at,
        published_at=resp.time.published_at,
        updated_at=resp.time.published_at,
        closed_at=resp.time.closed_at,
        moderated_at=resp.time.moderated_at,
        moderator_comments=resp.stats.moderator_comments,
        views=resp.stats.views,
        applications_count=resp.stats.applications_count,
    )
