from hypothesis import given
from hypothesis import strategies as st
from pkg.vacancy_api.vacancy_pb2 import Currency as PKGCurrency
from pkg.vacancy_api.vacancy_pb2 import RemoteType as PKGRemoteType
from pkg.vacancy_api.vacancy_pb2 import TimeType as PKGTimeType
from pkg.vacancy_api.vacancy_pb2 import VacancyInfo as PKGVacancyInfo

from src.api.grpc.dto.vacancy import vacancy_response_dto
from src.domain.schemas import VacancyResponseSchema
from src.domain.types.enums import Currency, RemoteType, TimeType
from tests.unit.dto.asserts import (
    assert_additional,
    assert_currencies,
    assert_main,
    assert_remote_and_time,
    assert_salary,
    assert_times,
)
from tests.unit.dto.vacancy_info_gen import FullVacancyInfo, full_vacancy_info


@given(responses=st.lists(full_vacancy_info(), min_size=5, max_size=5))
def test_vacancy_response_dto(responses: list[FullVacancyInfo]) -> None:
    for resp in responses:
        vacancy = _create_vacancy_response(resp)
        info = vacancy_response_dto(vacancy)

        assert_vacancy_response(vacancy, info)


def assert_vacancy_response(schema: VacancyResponseSchema, info: PKGVacancyInfo) -> None:
    assert info.vacancy_id == schema.vacancy_id
    assert info.description == schema.description

    assert_main(info, schema)
    assert_salary(info, schema)
    assert_currencies(info, schema)
    assert_salary(info, schema)

    assert info.status == schema.status.value

    assert_remote_and_time(info, schema)
    assert_times(info, schema)
    assert_additional(info, schema)

    assert info.views == schema.views
    assert info.applications_count == schema.applications_count

    assert info.author_name == schema.author_name

    if schema.moderator_comments is None:
        assert info.moderator_comments == ""
    else:
        assert info.moderator_comments == schema.moderator_comments


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
        updated_at=resp.time.updated_at,
        closed_at=resp.time.closed_at,
        moderated_at=resp.time.moderated_at,
        moderator_comments=resp.stats.moderator_comments,
        views=resp.stats.views,
        applications_count=resp.stats.applications_count,
    )
