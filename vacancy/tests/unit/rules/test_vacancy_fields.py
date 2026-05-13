from pytest import raises

from src.core.exc import ArgumentError
from src.domain.rules.vacancy_fields import validate_salaries


def test_validate_salaries(valid_salaries, invalid_salaries) -> None:
    for salaries in valid_salaries:
        salary_min, salary_max = salaries
        validate_salaries(salary_min, salary_max)

    for salaries in invalid_salaries:
        salary_min, salary_max = salaries

        with raises(ArgumentError):
            validate_salaries(salary_min, salary_max)
