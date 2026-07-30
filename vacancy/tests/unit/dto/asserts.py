from datetime import datetime
from typing import Any, Protocol

from google.protobuf.timestamp import Timestamp
from pkg.vacancy_api.vacancy_pb2 import CreateVacancyRequest
from pkg.vacancy_api.vacancy_pb2 import Currency as PKGCurrency
from pkg.vacancy_api.vacancy_pb2 import RemoteType as PKGRemoteType
from pkg.vacancy_api.vacancy_pb2 import TimeType as PKGTimeType

from src.domain.schemas import VacancyCreateSchema
from src.domain.types import Currency, Money, RemoteType, TimeType, Year


class VacancyMain(Protocol):
    title: str
    conditions: str
    requirements: str
    description: str
    tags: Any


class VacancySalary(Protocol):
    salary_min: Money | None
    salary_max: Money | None


class VacancyCurrency(Protocol):
    currency: Currency


class VacancyAdditional(Protocol):
    city: str | None
    metro: str | None
    is_city_valid: bool
    is_metro_valid: bool


class VacancyRemoteAndTime(Protocol):
    remote_type: RemoteType
    time_type: TimeType


class VacancyExp(Protocol):
    experience_min: Year | None
    experience_max: Year | None


class VacancyTime(Protocol):
    created_at: datetime
    updated_at: datetime | None
    published_at: datetime | None
    closed_at: datetime | None
    moderated_at: datetime | None


def assert_main(vacancy, schema: VacancyMain) -> None:
    assert vacancy.title == schema.title
    assert vacancy.conditions == schema.conditions
    assert vacancy.requirements == schema.requirements

    assert list(vacancy.tags) == schema.tags


def assert_salary(vacancy, schema: VacancySalary) -> None:
    if schema.salary_min is None:
        assert vacancy.salary_min == 0
    else:
        assert vacancy.salary_min == schema.salary_min

    if schema.salary_max is None:
        assert vacancy.salary_max == 0
    else:
        assert vacancy.salary_max == schema.salary_max


def assert_currencies(vacancy, schema: VacancyCurrency) -> None:
    currs = {
        PKGCurrency.EUR: Currency.EUR,
        PKGCurrency.USD: Currency.USD,
        PKGCurrency.RUB: Currency.RUB,
    }

    assert currs[vacancy.currency] == schema.currency


def assert_additional(vacancy, schema: VacancyAdditional) -> None:
    if schema.city is None:
        assert vacancy.city == ""
    else:
        assert vacancy.city == schema.city

    if schema.metro is None:
        assert vacancy.metro == ""
    else:
        assert vacancy.metro == schema.metro

    if hasattr(vacancy, "is_metro_valid"):
        assert vacancy.is_metro_valid is schema.is_metro_valid

    if hasattr(vacancy, "is_city_valid"):
        assert vacancy.is_city_valid is schema.is_city_valid


def assert_remote_and_time(vacancy, schema: VacancyRemoteAndTime) -> None:
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


def assert_exp(vacancy, schema: VacancyExp) -> None:
    if schema.experience_min is None:
        assert vacancy.experience_min == 0
    else:
        assert vacancy.experience_min == schema.experience_min

    if schema.experience_max is None:
        assert vacancy.experience_max == 0
    else:
        assert vacancy.experience_max == schema.experience_max


def assert_author(vacancy: CreateVacancyRequest, schema: VacancyCreateSchema) -> None:
    assert vacancy.user_info.company_name == schema.author_name
    assert vacancy.user_info.user_id == str(schema.author_id)


def assert_times(vacancy, schema: VacancyTime) -> None:
    def assert_time(field_name: str, eq: datetime | None) -> None:
        attr = getattr(vacancy, field_name)
        if eq is None:
            assert attr == Timestamp()
        else:
            assert attr.ToDatetime() == eq

    assert_time("created_at", schema.created_at)
    assert_time("moderated_time", schema.moderated_at)
    assert_time("updated_at", schema.updated_at)
    assert_time("published_at", schema.published_at)
    assert_time("closed_at", schema.closed_at)
