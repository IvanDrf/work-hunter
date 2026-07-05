from re import escape
from string import Template

from hypothesis import given
from hypothesis import strategies as st
from pydantic import BaseModel
from pytest import raises

from src.core.exc import ArgumentError
from src.domain.schemas.mixins import MAX_YEAR, MIN_YEAR, ExperienceValidatorMixin
from src.domain.types import UNSET_VALUE, Money, UnsetValue


class ExperienceTestClass(BaseModel, ExperienceValidatorMixin):
    experience_min: Money | None | UnsetValue
    experience_max: Money | None | UnsetValue


@given(
    negative_exps=st.lists(st.integers(-MAX_YEAR, MIN_YEAR - 1), min_size=10, max_size=30),
    valid_exps=st.lists(st.integers(MIN_YEAR, MAX_YEAR), min_size=10, max_size=30),
    invalid_big_exps=st.lists(st.integers(MAX_YEAR + 1, MAX_YEAR * 2), min_size=10, max_size=30),
)
def test_ExperienceValidatorMixin(negative_exps: list[int], valid_exps: list[int], invalid_big_exps: list[int]) -> None:
    for negative_experience in negative_exps:
        _test_unbound_experience(negative_experience)

    for valid_experience in valid_exps:
        _test_valid_Experience(valid_experience)

    for unbound_exp in invalid_big_exps:
        _test_unbound_experience(unbound_exp)

    INVALID_EXPERIENCE_VALUES = Template(
        "maximum experience - $experience_max must be greater than minimum experience - $experience_min"
    )
    valid_exps.sort()

    for i in range(len(valid_exps) - 1):
        experience_min, experience_max = valid_exps[i], valid_exps[i + 1]

        e = ExperienceTestClass(experience_min=experience_min, experience_max=experience_max)
        assert e.experience_min == experience_min
        assert e.experience_max == experience_max

        if experience_min != experience_max:
            with raises(
                ArgumentError,
                match=INVALID_EXPERIENCE_VALUES.substitute(experience_min=experience_max, experience_max=experience_min),
            ):
                ExperienceTestClass(experience_min=experience_max, experience_max=experience_min)


def _test_unbound_experience(experience: int) -> None:
    INVALID_EXPERIENCE_RANGE = Template(f"experience value must be in range ({MIN_YEAR}, {MAX_YEAR}), but given=$experience")

    # valid values for other experience: max or min
    # example: experience_min = -200, experience_max = None
    for valid_experience in (None, UNSET_VALUE, 0):
        error_message = escape(INVALID_EXPERIENCE_RANGE.substitute(experience=experience))

        with raises(ArgumentError, match=error_message):
            ExperienceTestClass(experience_min=experience, experience_max=valid_experience)

        e = ExperienceTestClass(experience_min=None, experience_max=None)
        with raises(ArgumentError, match=error_message):
            e.experience_min = experience

        with raises(ArgumentError, match=error_message):
            ExperienceTestClass(experience_min=valid_experience, experience_max=experience)

        e = ExperienceTestClass(experience_min=None, experience_max=None)
        with raises(ArgumentError, match=error_message):
            e.experience_max = experience


def _test_valid_Experience(valid_experience: int) -> None:
    """
    Test valid Experience with default values
        example: Experience_min = 50, Experience_max = None/UNSET_VALUE,
        example: Experience_min = 0/UNSET_VALUE/None, Experience_max = 50
    """
    for valid_experience_def in (None, UNSET_VALUE, 0):
        # check values for Experience_min
        # if valid_Experience == 0 there is ony one valid value for Experience_min it's zero
        e = ExperienceTestClass(
            experience_min=valid_experience, experience_max=valid_experience_def if valid_experience_def != 0 else None
        )
        assert e.experience_min == valid_experience
        if valid_experience_def != 0:
            assert e.experience_max == valid_experience_def
        else:
            assert e.experience_max is None

        e = ExperienceTestClass(experience_min=None, experience_max=None)
        assert e.experience_min is None
        assert e.experience_max is None

        e.experience_max = valid_experience_def
        assert e.experience_min is None
        assert e.experience_max == valid_experience_def

        # check values for Experience_max
        e = ExperienceTestClass(experience_min=valid_experience_def, experience_max=valid_experience)
        assert e.experience_min == valid_experience_def
        assert e.experience_max == valid_experience

        e = ExperienceTestClass(experience_min=None, experience_max=None)
        assert e.experience_min is None
        assert e.experience_max is None

        e.experience_min = valid_experience_def
        assert e.experience_min == valid_experience_def
        assert e.experience_max is None
