from dataclasses import dataclass
from uuid import UUID

from hypothesis import given
from hypothesis import strategies as st
from pkg.common.common_pb2 import FullUserInfo as PKGFullUserInfo
from pkg.common.common_pb2 import UserRole as PKGUserRole
from pkg.vacancy_api.vacancy_pb2 import CreateVacancyRequest
from pkg.vacancy_api.vacancy_pb2 import Currency as PKGCurrency
from pkg.vacancy_api.vacancy_pb2 import RemoteType as PKGRemoteType
from pkg.vacancy_api.vacancy_pb2 import TimeType as PKGTimeType

from src.api.dto.vacancy import vacancy_create_dto
from src.domain.schemas import VacancyCreateSchema
from src.domain.schemas.mixins import MAX_MONEY, MAX_YEAR, MIN_MONEY, MIN_YEAR
from src.domain.schemas.vacancy import (
    MAX_AUTHOR_NAME_LENGTH,
    MAX_CITY_LENGTH,
    MAX_CONDITIONS_LENGTH,
    MAX_DESCRIPTION_LENGTH,
    MAX_METRO_LENGTH,
    MAX_REQUIREMENTS_LENGTH,
    MAX_TAGS_AMOUNT,
    MAX_TITLE_LENGTH,
    MIN_TAGS_AMOUNT,
    MIN_TITLE_LENGTH,
)
from src.domain.types.enums import Currency, RemoteType, TimeType
from src.domain.types.types import Money, Year


@dataclass(frozen=True)
class VacancyMain:
    title: str
    description: str
    requirements: str
    conditions: str

    tags: list[str]


@dataclass(frozen=True)
class VacancySalary:
    salary_min: Money
    salary_max: Money
    currency: PKGCurrency


@dataclass(frozen=True)
class VacancyAdditional:
    city: str | None
    metro: str | None
    remote_type: PKGRemoteType
    time_type: PKGTimeType


@dataclass(frozen=True)
class VacancyExp:
    experience_min: Year
    experience_max: Year


@dataclass(frozen=True)
class VacancyAuthor:
    author_id: UUID
    author_name: str
    role: PKGUserRole
    verificated: bool


@dataclass(frozen=True)
class VacancyInfo:
    main: VacancyMain
    salary: VacancySalary
    additional: VacancyAdditional
    exp: VacancyExp
    author: VacancyAuthor


@st.composite
def vacancy_main_valid(draw) -> VacancyMain:
    title = draw(st.text(min_size=MIN_TITLE_LENGTH, max_size=MAX_TITLE_LENGTH))
    description = draw(st.text(max_size=MAX_DESCRIPTION_LENGTH))
    requirements = draw(st.text(max_size=MAX_REQUIREMENTS_LENGTH))
    conditions = draw(st.text(max_size=MAX_CONDITIONS_LENGTH))

    tags = draw(
        st.lists(
            st.sampled_from(["go", "cpp", "worker", "python", "rabbitmq", "ozon"]),
            min_size=MIN_TAGS_AMOUNT,
            max_size=MAX_TAGS_AMOUNT,
        )
    )

    return VacancyMain(title=title, description=description, requirements=requirements, conditions=conditions, tags=tags)


@st.composite
def vacancy_salary_valid(draw) -> VacancySalary:
    salary_min = draw(st.integers(min_value=MIN_MONEY, max_value=MAX_MONEY // 2))
    salary_max = draw(st.integers(min_value=MAX_MONEY // 2, max_value=MAX_MONEY))
    currency = draw(st.sampled_from([PKGCurrency.EUR, PKGCurrency.USD, PKGCurrency.RUB]))

    return VacancySalary(salary_min=salary_min, salary_max=salary_max, currency=currency)


@st.composite
def vacancy_additional_valid(draw) -> VacancyAdditional:
    city = draw(st.text(max_size=MAX_CITY_LENGTH))
    metro = draw(st.text(max_size=MAX_METRO_LENGTH))

    remote = draw(st.sampled_from([PKGRemoteType.ANY, PKGRemoteType.HYBRID, PKGRemoteType.OFFICE, PKGRemoteType.REMOTE]))
    time_type = draw(st.sampled_from([PKGTimeType.FULL, PKGTimeType.PART, PKGTimeType.PART]))

    return VacancyAdditional(city=city, metro=metro, remote_type=remote, time_type=time_type)


@st.composite
def vacancy_exp_valid(draw) -> VacancyExp:
    exp_min = draw(st.integers(min_value=MIN_YEAR, max_value=MAX_YEAR // 2))
    exp_max = draw(st.integers(min_value=MAX_YEAR // 2, max_value=MAX_YEAR))

    return VacancyExp(experience_min=exp_min, experience_max=exp_max)


@st.composite
def vacancy_author_valid(draw) -> VacancyAuthor:
    author_id = draw(st.uuids(version=4))
    author_name = draw(st.text(max_size=MAX_AUTHOR_NAME_LENGTH))
    role = draw(st.sampled_from([PKGUserRole.ADMIN, PKGUserRole.EMPLOYEE, PKGUserRole.EMPLOYER]))
    verificated = draw(st.sampled_from([False, True]))

    return VacancyAuthor(author_id=author_id, author_name=author_name, role=role, verificated=verificated)


@st.composite
def vacancy_info(draw) -> VacancyInfo:
    return VacancyInfo(
        main=draw(vacancy_main_valid()),
        salary=draw(vacancy_salary_valid()),
        additional=draw(vacancy_additional_valid()),
        exp=draw(vacancy_exp_valid()),
        author=draw(vacancy_author_valid()),
    )


@given(requests=st.lists(vacancy_info(), min_size=5, max_size=10))
def test_vacancy_create_dto(requests: list[VacancyInfo]) -> None:
    for req in requests:
        r = _create_request(req)

        schema = vacancy_create_dto(r)
        assert_main(r, schema)
        assert_salary(r, schema)
        assert_additional(r, schema)
        assert_exp(r, schema)
        assert_author(r, schema)


@given(requests=st.lists(vacancy_info(), min_size=5, max_size=10))
def test_vacancy_create_dto_no_optional_fields(requests: list[VacancyInfo]) -> None:
    for req in requests:
        r = _create_request_no_optional_fields(req)

        schema = vacancy_create_dto(r)
        assert schema.title == r.title
        assert schema.description == "No description"
        assert schema.requirements == r.requirements
        assert schema.conditions == r.conditions

        assert schema.salary_min is None
        assert schema.salary_max is None
        assert schema.currency == Currency.RUB

        assert schema.city is None
        assert schema.metro is None
        assert_remote_and_time(r, schema)

        assert schema.experience_min is None
        assert schema.experience_max is None
        assert schema.tags == list(r.tags)


def assert_main(vacancy: CreateVacancyRequest, schema: VacancyCreateSchema) -> None:
    assert vacancy.title == schema.title
    assert vacancy.conditions == schema.conditions
    assert vacancy.requirements == schema.requirements
    assert vacancy.tags == schema.tags


def assert_salary(vacancy: CreateVacancyRequest, schema: VacancyCreateSchema) -> None:
    assert vacancy.salary_min == schema.salary_min
    assert vacancy.salary_max == schema.salary_max

    currs = {
        PKGCurrency.EUR: Currency.EUR,
        PKGCurrency.USD: Currency.USD,
        PKGCurrency.RUB: Currency.RUB,
    }

    assert currs[vacancy.currency] == schema.currency


def assert_additional(vacancy: CreateVacancyRequest, schema: VacancyCreateSchema) -> None:
    assert vacancy.city == schema.city
    assert vacancy.metro == schema.metro

    assert_remote_and_time(vacancy, schema)


def assert_remote_and_time(vacancy: CreateVacancyRequest, schema: VacancyCreateSchema) -> None:
    remotes = {
        PKGRemoteType.ANY: RemoteType.ANY,
        PKGRemoteType.HYBRID: RemoteType.HYBRID,
        PKGRemoteType.OFFICE: RemoteType.OFFICE,
        PKGRemoteType.REMOTE: RemoteType.REMOTE,
    }

    assert remotes[vacancy.remote_type] == schema.remote_type

    time_types = {
        PKGTimeType.FULL: TimeType.FULL,
        PKGTimeType.PART: TimeType.PART,
    }

    assert time_types[vacancy.time_type] == schema.time_type


def assert_exp(vacancy: CreateVacancyRequest, schema: VacancyCreateSchema) -> None:
    assert vacancy.experience_min == schema.experience_min
    assert vacancy.experience_max == schema.experience_max


def assert_author(vacancy: CreateVacancyRequest, schema: VacancyCreateSchema) -> None:
    assert vacancy.user_info.username == schema.author_name
    assert vacancy.user_info.user_id == str(schema.author_id)


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
