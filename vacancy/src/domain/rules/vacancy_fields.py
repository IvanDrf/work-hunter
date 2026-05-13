from src.core.exc import ArgumentError
from src.domain.types.types import UNSET_VALUE, Money, UnsetValue


def validate_salaries(salary_min: Money | None | UnsetValue, salary_max: Money | None | UnsetValue) -> None:
    if salary_min is not None and salary_min is not UNSET_VALUE:
        validate_salary(salary_min)

    if salary_max is not None and salary_max is not UNSET_VALUE:
        validate_salary(salary_max)

    if type(salary_min) is Money and type(salary_max) is Money and salary_min > salary_max:
        raise ArgumentError(f"maximum salary - {salary_max} must be greater than minimum salary - {salary_min}")


def validate_salary(salary: Money | UnsetValue) -> None:
    if salary < 0:
        raise ArgumentError("salary must be a non negative value")
