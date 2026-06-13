from hypothesis import given
from hypothesis import strategies as st
from pkg.common.common_pb2 import FullUserInfo as PKGFullUserInfo
from pkg.vacancy_api.vacancy_pb2 import CreateVacancyRequest

from src.api.dto.vacancy import vacancy_create_dto
from src.domain.types.enums import Currency
from tests.unit.dto.asserts import (
    assert_additional,
    assert_author,
    assert_currencies,
    assert_exp,
    assert_main,
    assert_remote_and_time,
    assert_salary,
)
from tests.unit.dto.vacancy_info_gen import VacancyInfo, vacancy_info


@given(requests=st.lists(vacancy_info(), min_size=5, max_size=10))
def test_vacancy_create_dto(requests: list[VacancyInfo]) -> None:
    for req in requests:
        r = _create_request(req)

        schema = vacancy_create_dto(r)
        assert_main(r, schema)
        assert_salary(r, schema)
        assert_currencies(r, schema)
        assert_additional(r, schema)
        assert_remote_and_time(r, schema)
        assert_exp(r, schema)
        assert_author(r, schema)


@given(requests=st.lists(vacancy_info(), min_size=5, max_size=10))
def test_vacancy_create_dto_no_optional_fields(requests: list[VacancyInfo]) -> None:
    for req in requests:
        r = _create_request_no_optional_fields(req)

        schema = vacancy_create_dto(r)
        assert_main(r, schema)

        assert_salary(r, schema)
        assert schema.currency == Currency.RUB

        assert_additional(r, schema)
        assert_remote_and_time(r, schema)
        assert_exp(r, schema)


def _create_request(req: VacancyInfo) -> CreateVacancyRequest:
    return CreateVacancyRequest(
        title=req.main.title,
        description=req.main.description,
        conditions=req.main.conditions,
        requirements=req.main.requirements,
        salary_min=req.salary.salary_min,
        salary_max=req.salary.salary_max,
        currency=req.salary.currency,
        city=req.additional.city,
        metro=req.additional.metro,
        remote_type=req.additional.remote_type,
        time_type=req.additional.time_type,
        experience_min=req.exp.experience_min,
        experience_max=req.exp.experience_min,
        tags=req.main.tags,
        user_info=PKGFullUserInfo(
            role=req.author.role,
            user_id=str(req.author.author_id),
            verificated=req.author.verificated,
            username=req.author.author_name,
        ),
    )


def _create_request_no_optional_fields(req: VacancyInfo) -> CreateVacancyRequest:
    return CreateVacancyRequest(
        title=req.main.title,
        requirements=req.main.requirements,
        conditions=req.main.conditions,
        remote_type=req.additional.remote_type,
        time_type=req.additional.time_type,
        tags=req.main.tags,
        user_info=PKGFullUserInfo(
            role=req.author.role,
            user_id=str(req.author.author_id),
            verificated=req.author.verificated,
            username=req.author.author_name,
        ),
    )
