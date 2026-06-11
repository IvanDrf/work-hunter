from src.domain.rules.user import is_user_admin, is_user_employer


def test_is_user_admin(admin_user_info, employee_user_info, employer_user_info, unspecified_user_info) -> None:
    for admin, employee, employer, unspec in zip(admin_user_info, employee_user_info, employer_user_info, unspecified_user_info):
        assert is_user_admin(admin) is True

        assert is_user_admin(employee) is False
        assert is_user_admin(employer) is False
        assert is_user_admin(unspec) is False


def test_is_user_employer(admin_user_info, employee_user_info, employer_user_info, unspecified_user_info) -> None:
    for admin, employee, employer, unspec in zip(admin_user_info, employee_user_info, employer_user_info, unspecified_user_info):
        assert is_user_employer(employer) is True

        assert is_user_employer(admin) is False
        assert is_user_employer(employee) is False
        assert is_user_employer(unspec) is False
