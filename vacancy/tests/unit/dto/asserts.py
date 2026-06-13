from typing import Any, Protocol

from pkg.vacancy_api.vacancy_pb2 import CreateVacancyRequest
from pkg.vacancy_api.vacancy_pb2 import Currency as PKGCurrency
from pkg.vacancy_api.vacancy_pb2 import RemoteType as PKGRemoteType
from pkg.vacancy_api.vacancy_pb2 import TimeType as PKGTimeType

from src.domain.schemas import VacancyCreateSchema
from src.domain.types.enums import Currency, RemoteType, TimeType
from src.domain.types.types import Money, Year


class VacancyMain(Protocol):
    title: str
    conditions: str
    requirements: str
    tags: Any


class VacancySalary(Protocol):
    salary_min: Money | None
    salary_max: Money | None


def assert_main(vacancy: VacancyMain, schema: VacancyMain) -> None:
    assert vacancy.title == schema.title
    assert vacancy.conditions == schema.conditions
    assert vacancy.requirements == schema.requirements
    assert list(vacancy.tags) == schema.tags


def assert_salary(vacancy: VacancySalary, schema: VacancySalary) -> None:
    assert vacancy.salary_min == schema.salary_min
    assert vacancy.salary_max == schema.salary_max


def assert_currencies(vacancy: CreateVacancyRequest, schema: VacancyCreateSchema) -> None:
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
