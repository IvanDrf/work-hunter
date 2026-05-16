from pytest import mark

from src.domain.rules.vacancy import is_vacancy_id_valid


@mark.parametrize("vacancy_id", [1, -5, 2, 19, -3, 0, 23, -39, 123])
def test_is_vacancy_id_valid(vacancy_id) -> None:
    assert is_vacancy_id_valid(vacancy_id) is (vacancy_id > 0)
