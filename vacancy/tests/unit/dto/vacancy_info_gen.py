from dataclasses import dataclass
from datetime import datetime
from uuid import UUID

from hypothesis import strategies as st
from pkg.common.common_pb2 import UserRole as PKGUserRole
from pkg.vacancy_api.vacancy_pb2 import Currency as PKGCurrency
from pkg.vacancy_api.vacancy_pb2 import RemoteType as PKGRemoteType
from pkg.vacancy_api.vacancy_pb2 import TimeType as PKGTimeType

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
from src.domain.types.enums import VacancyStatus
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
class VacancyTime:
    created_at: datetime
    updated_at: datetime
    published_at: datetime
    closed_at: datetime
    moderated_at: datetime


@dataclass(frozen=True)
class VacancyStats:
    status: VacancyStatus
    moderator_comments: str | None
    views: int
    applications_count: int


@dataclass(frozen=True)
class VacancyInfo:
    main: VacancyMain
    salary: VacancySalary
    additional: VacancyAdditional
    exp: VacancyExp
    author: VacancyAuthor


@dataclass(frozen=True)
class FullVacancyInfo(VacancyInfo):
    vacancy_id: int
    time: VacancyTime
    stats: VacancyStats


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
    city = draw(st.text(min_size=1, max_size=MAX_CITY_LENGTH))
    metro = draw(st.text(min_size=1, max_size=MAX_METRO_LENGTH))

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
    author_name = draw(st.text(min_size=1, max_size=MAX_AUTHOR_NAME_LENGTH))
    role = draw(st.sampled_from([PKGUserRole.ADMIN, PKGUserRole.EMPLOYEE, PKGUserRole.EMPLOYER]))
    verificated = draw(st.one_of(st.booleans()))

    return VacancyAuthor(author_id=author_id, author_name=author_name, role=role, verificated=verificated)


@st.composite
def vacancy_time_valid(draw) -> VacancyTime:
    created_at = draw(st.datetimes())
    updated_at = draw(st.one_of(st.datetimes(), st.none()))
    published_at = draw(st.one_of(st.datetimes(), st.none()))
    closed_at = draw(st.one_of(st.datetimes(), st.none()))
    moderated_at = draw(st.one_of(st.datetimes(), st.none()))

    return VacancyTime(
        created_at=created_at,
        updated_at=updated_at,
        published_at=published_at,
        closed_at=closed_at,
        moderated_at=moderated_at,
    )


@st.composite
def vacancy_stats_valid(draw) -> VacancyStats:
    status = draw(st.sampled_from(VacancyStatus))
    moderator_comments = draw(st.one_of(st.none(), st.text(min_size=0, max_size=200)))
    views = draw(st.integers(min_value=0, max_value=2000))
    applications_count = draw(st.integers(min_value=0, max_value=2000))

    return VacancyStats(status=status, moderator_comments=moderator_comments, views=views, applications_count=applications_count)


@st.composite
def vacancy_info(draw) -> VacancyInfo:
    return VacancyInfo(
        main=draw(vacancy_main_valid()),
        salary=draw(vacancy_salary_valid()),
        additional=draw(vacancy_additional_valid()),
        exp=draw(vacancy_exp_valid()),
        author=draw(vacancy_author_valid()),
    )


@st.composite
def full_vacancy_info(draw) -> FullVacancyInfo:
    return FullVacancyInfo(
        vacancy_id=draw(st.integers(min_value=0, max_value=90000)),
        main=draw(vacancy_main_valid()),
        salary=draw(vacancy_salary_valid()),
        additional=draw(vacancy_additional_valid()),
        exp=draw(vacancy_exp_valid()),
        author=draw(vacancy_author_valid()),
        time=draw(vacancy_time_valid()),
        stats=draw(vacancy_stats_valid()),
    )
