from pytest import mark

from src.domain.rules.vacancy import has_right_to_vacancy, is_vacancy_id_valid


@mark.parametrize("vacancy_id", [1, -5, 2, 19, -3, 0, 23, -39, 123])
def test_is_vacancy_id_valid(vacancy_id) -> None:
    assert is_vacancy_id_valid(vacancy_id) is (vacancy_id > 0)


def test_has_right_to_vacancy(vacancies) -> None:
    for vacancy, user_info, res in vacancies:
        assert has_right_to_vacancy(vacancy, user_info) is res
