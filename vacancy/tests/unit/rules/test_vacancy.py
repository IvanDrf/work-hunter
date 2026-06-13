from hypothesis import given
from hypothesis import strategies as st

from src.domain.rules.vacancy import has_right_to_vacancy, is_vacancy_id_valid


@given(vacancy_id=st.integers(min_value=-100, max_value=100))
def test_is_vacancy_id_valid(vacancy_id) -> None:
    assert is_vacancy_id_valid(vacancy_id) is (vacancy_id > 0)


def test_has_right_to_vacancy(vacancies) -> None:
    for vacancy, user_info, res in vacancies:
        assert has_right_to_vacancy(vacancy, user_info) is res
