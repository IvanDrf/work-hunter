from hypothesis import given
from hypothesis import strategies as st

from src.domain.types.types import UNSET_VALUE, UnsetValue


@given(nums=st.lists(st.integers(-20_000, 20_000), min_size=10, max_size=30))
def test_unsetValue(nums: list[int]) -> None:
    for num in nums:
        unset = UnsetValue()

        assert unset is UNSET_VALUE  # singleton check
        assert (unset < num) is True
        assert (unset > num) is False
