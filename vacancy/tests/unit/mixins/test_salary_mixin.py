from re import escape
from string import Template

from hypothesis import given
from hypothesis import strategies as st
from pydantic import BaseModel
from pytest import raises

from src.core.exc import ArgumentError
from src.domain.schemas.mixins import MAX_MONEY, MIN_MONEY, SalaryValidatorMixin
from src.domain.types.types import UNSET_VALUE, Money, UnsetValue


class SalaryTestClass(BaseModel, SalaryValidatorMixin):
    salary_min: Money | None | UnsetValue
    salary_max: Money | None | UnsetValue


@given(
    negative_salaries=st.lists(st.integers(-MAX_MONEY, MIN_MONEY - 1), min_size=10, max_size=30),
    valid_salaries=st.lists(st.integers(MIN_MONEY, MAX_MONEY), min_size=10, max_size=30),
    invalid_big_salaries=st.lists(st.integers(MAX_MONEY + 1, MAX_MONEY * 2), min_size=10, max_size=30),
)
def test_SalaryValidatorMixin(negative_salaries: list[int], valid_salaries: list[int], invalid_big_salaries: list[int]) -> None:
    for negative_salary in negative_salaries:
        _test_unbound_salary(negative_salary)

    for valid_salary in valid_salaries:
        _test_valid_salary(valid_salary)

    for invalid_big_salary in invalid_big_salaries:
        _test_unbound_salary(invalid_big_salary)

    INVALID_SALARY_VALUES = Template("maximum salary - $salary_max must be greater than minimum salary - $salary_min")
    valid_salaries.sort()

    for i in range(len(valid_salaries) - 1):
        salary_min, salary_max = valid_salaries[i], valid_salaries[i + 1]

        s = SalaryTestClass(salary_min=salary_min, salary_max=salary_max)
        assert s.salary_min == salary_min
        assert s.salary_max == salary_max

        if salary_min != salary_max:
            with raises(ArgumentError, match=INVALID_SALARY_VALUES.substitute(salary_min=salary_max, salary_max=salary_min)):
                SalaryTestClass(salary_min=salary_max, salary_max=salary_min)


def _test_unbound_salary(salary: int) -> None:
    INVALID_SALARY_RANGE = Template(f"salary value must be in range ({MIN_MONEY}, {MAX_MONEY}), but given=$salary")

    # valid values for other salary: max or min
    # example: salary_min = -200, salary_max = None
    for valid_salary in (None, UNSET_VALUE, 0):
        error_message = escape(INVALID_SALARY_RANGE.substitute(salary=salary))

        with raises(ArgumentError, match=error_message):
            SalaryTestClass(salary_min=salary, salary_max=valid_salary)

        s = SalaryTestClass(salary_min=None, salary_max=None)
        with raises(ArgumentError, match=error_message):
            s.salary_min = salary

        with raises(ArgumentError, match=error_message):
            SalaryTestClass(salary_min=valid_salary, salary_max=salary)

        s = SalaryTestClass(salary_min=None, salary_max=None)
        with raises(ArgumentError, match=error_message):
            s.salary_max = salary


def _test_valid_salary(valid_salary: int) -> None:
    """
    Test valid salary with default values
        example: salary_min = 50, salary_max = None/UNSET_VALUE,
        example: salary_min = 0/UNSET_VALUE/None, salary_max = 50
    """
    for valid_salary_def in (None, UNSET_VALUE, 0):
        # check values for salary_min
        # if valid_salary == 0 there is ony one valid value for salary_min it's zero
        s = SalaryTestClass(salary_min=valid_salary, salary_max=valid_salary_def if valid_salary_def != 0 else None)
        assert s.salary_min == valid_salary
        if valid_salary_def != 0:
            assert s.salary_max == valid_salary_def
        else:
            assert s.salary_max is None

        s = SalaryTestClass(salary_min=None, salary_max=None)
        assert s.salary_min is None
        assert s.salary_max is None

        s.salary_max = valid_salary_def
        assert s.salary_min is None
        assert s.salary_max == valid_salary_def

        # check values for salary_max
        s = SalaryTestClass(salary_min=valid_salary_def, salary_max=valid_salary)
        assert s.salary_min == valid_salary_def
        assert s.salary_max == valid_salary

        s = SalaryTestClass(salary_min=None, salary_max=None)
        assert s.salary_min is None
        assert s.salary_max is None

        s.salary_min = valid_salary_def
        assert s.salary_min == valid_salary_def
        assert s.salary_max is None
