from datetime import datetime
from typing import Final
from uuid import UUID

from pydantic import BaseModel, Field

from src.domain.schemas.mixins import ExperienceValidatorMixin, SalaryValidatorMixin
from src.domain.types.enums import Currency, RemoteType, TimeType, VacancyStatus
from src.domain.types.types import UNSET_VALUE, Money, UnsetValue, Year

MIN_TITLE_LENGTH: Final[int] = 5
MAX_TITLE_LENGTH: Final[int] = 150

MAX_DESCRIPTION_LENGTH: Final[int] = 20_000
MAX_REQUIREMENTS_LENGTH: Final[int] = 20_000
MAX_CONDITIONS_LENGTH: Final[int] = 20_000


class VacancySchema(BaseModel, SalaryValidatorMixin, ExperienceValidatorMixin):
    title: str = Field(min_length=MIN_TITLE_LENGTH, max_length=MAX_TITLE_LENGTH)
    description: str = Field(default="No description", max_length=MAX_DESCRIPTION_LENGTH)

    requirements: str = Field(max_length=MAX_REQUIREMENTS_LENGTH)
    conditions: str = Field(max_length=MAX_CONDITIONS_LENGTH)

    salary_min: Money | None = None
    salary_max: Money | None = None

    currency: Currency = Currency.RUB

    city: str | None = None
    metro: str | None = None

    remote_type: RemoteType = Field(default=RemoteType.ANY)
    time_type: TimeType = Field(default=TimeType.FULL)

    experience_min: Year | None = None
    experience_max: Year | None = None

    tags: list[str] = []


class VacancyCreateSchema(VacancySchema):
    author_id: UUID
    author_name: str


class VacancyUpdateSchema(BaseModel, SalaryValidatorMixin, ExperienceValidatorMixin):
    vacancy_id: int = Field(ge=0)
    title: str | UnsetValue = UNSET_VALUE
    description: str | UnsetValue = UNSET_VALUE

    requirements: str | UnsetValue = UNSET_VALUE
    conditions: str | UnsetValue = UNSET_VALUE

    salary_min: Money | None | UnsetValue = UNSET_VALUE
    salary_max: Money | None | UnsetValue = UNSET_VALUE

    currency: Currency | UnsetValue = UNSET_VALUE

    city: str | None | UnsetValue = UNSET_VALUE
    metro: str | None | UnsetValue = UNSET_VALUE

    remote_type: RemoteType | UnsetValue = UNSET_VALUE
    time_type: TimeType | UnsetValue = UNSET_VALUE

    experience_min: Year | None | UnsetValue = UNSET_VALUE
    experience_max: Year | None | UnsetValue = UNSET_VALUE

    tags: list[str] | UnsetValue = UNSET_VALUE


class VacancyResponseSchema(VacancySchema):
    vacancy_id: int = Field(ge=0)
    author_name: str
    author_id: UUID

    status: VacancyStatus = VacancyStatus.MODERATING

    views: int = Field(default=0, ge=0)
    applications_count: int = Field(default=0, ge=0)

    created_at: datetime
    updated_at: datetime | None = None
    published_at: datetime | None = None
    closed_at: datetime | None = None

    moderator_comments: str | None = None
    moderated_at: datetime | None = None
