from pydantic import ValidationInfo, field_validator

from src.core.exc import ArgumentError
from src.domain.types.types import UNSET_VALUE, Money, Year


class SalaryValidatorMixin:
    @field_validator("salary_min")
    def validate_salary_min(cls, value: Money, info: ValidationInfo) -> Money:
        salary_max = info.data.get("salary_max")
        if salary_max is not None and salary_max is not UNSET_VALUE and value > salary_max:
            raise ArgumentError(f"maximum salary - {salary_max} must be greater than minimum salary - {value}")

        return value

    @field_validator("salary_max")
    def validate_salary_max(cls, value: Money, info: ValidationInfo) -> Money:
        salary_min = info.data.get("salary_min")
        if salary_min is not None and salary_min is not UNSET_VALUE and value < salary_min:
            raise ArgumentError(f"maximum salary - {value} must be greater than minimum salary - {salary_min}")

        return value


class ExperienceValidatorMixin:
    @field_validator("experience_min")
    def validate_experience_min(cls, value: Year, info: ValidationInfo) -> Year:
        experience_max = info.data.get("experience_max")
        if experience_max is not None and experience_max is not UNSET_VALUE and value > experience_max:
            raise ArgumentError(f"maximum experience - {experience_max} must be greater than minimum experience - {value}")

        return value

    @field_validator("experience_max")
    def validate_experience_max(cls, value: Year, info: ValidationInfo) -> Year:
        experience_min = info.data.get("experience_min")
        if experience_min is not None and experience_min is not UNSET_VALUE and value < experience_min:
            raise ArgumentError(f"maximum experience - {value} must be greater than minimum experience - {experience_min}")

        return value
