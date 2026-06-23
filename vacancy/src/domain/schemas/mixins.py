from typing import Final

from pydantic import ValidationInfo, field_validator

from src.core.exc import ArgumentError
from src.domain.types.types import UNSET_VALUE, Money, UnsetValue, Year

MIN_MONEY: Final[Money] = 0
MAX_MONEY: Final[Money] = 1_000_000_000

MIN_YEAR: Final[Year] = 0
MAX_YEAR: Final[Year] = 1_000_000_000


class SalaryValidatorMixin:
    model_config = {"validate_assignment": True}

    @field_validator("salary_min")
    def validate_salary_min(cls, value: Money | UnsetValue | None, info: ValidationInfo) -> Money | UnsetValue | None:
        if value is None or isinstance(value, UnsetValue):
            return value

        if not (MIN_MONEY <= value <= MAX_MONEY):
            raise ArgumentError(f"salary value must be in range ({MIN_MONEY}, {MAX_MONEY}), but given={value}")

        salary_max = info.data.get("salary_max")
        if salary_max is not None and salary_max is not UNSET_VALUE and value > salary_max:
            raise ArgumentError(f"maximum salary - {salary_max} must be greater than minimum salary - {value}")

        return value

    @field_validator("salary_max")
    def validate_salary_max(cls, value: Money | UnsetValue | None, info: ValidationInfo) -> Money | UnsetValue | None:
        if value is None or isinstance(value, UnsetValue):
            return value

        if not (MIN_MONEY <= value <= MAX_MONEY):
            raise ArgumentError(f"salary value must be in range ({MIN_MONEY}, {MAX_MONEY}), but given={value}")

        salary_min = info.data.get("salary_min")
        if salary_min is not None and salary_min is not UNSET_VALUE and value < salary_min:
            raise ArgumentError(f"maximum salary - {value} must be greater than minimum salary - {salary_min}")

        return value


class ExperienceValidatorMixin:
    model_config = {"validate_assignment": True}

    @field_validator("experience_min")
    def validate_experience_min(cls, value: Year | UnsetValue | None, info: ValidationInfo) -> Year | UnsetValue | None:
        if value is None or isinstance(value, UnsetValue):
            return value

        if not (MIN_YEAR <= value <= MAX_YEAR):
            raise ArgumentError(f"experience value must be in range ({MIN_YEAR}, {MAX_YEAR}), but given={value}")

        experience_max = info.data.get("experience_max")
        if experience_max is not None and experience_max is not UNSET_VALUE and value > experience_max:
            raise ArgumentError(f"maximum experience - {experience_max} must be greater than minimum experience - {value}")

        return value

    @field_validator("experience_max")
    def validate_experience_max(cls, value: Year | UnsetValue | None, info: ValidationInfo) -> Year | UnsetValue | None:
        if value is None or isinstance(value, UnsetValue):
            return value

        if not (MIN_YEAR <= value <= MAX_YEAR):
            raise ArgumentError(f"experience value must be in range ({MIN_YEAR}, {MAX_YEAR}), but given={value}")

        experience_min = info.data.get("experience_min")
        if experience_min is not None and experience_min is not UNSET_VALUE and value < experience_min:
            raise ArgumentError(f"maximum experience - {value} must be greater than minimum experience - {experience_min}")

        return value
