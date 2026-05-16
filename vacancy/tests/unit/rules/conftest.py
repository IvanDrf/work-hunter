from uuid import UUID, uuid4

from pytest import fixture

from src.domain.models import VacancyORM
from src.domain.schemas.user import UserInfo
from src.domain.types.enums import UserRole, VacancyStatus


@fixture(scope="package")
def user_id() -> UUID:
    return UUID("3e0baeb1-a57d-42b2-bb05-a10cc3e5be57")


@fixture(scope="function")
def admin_user_info(user_id) -> tuple[UserInfo, ...]:
    return (
        UserInfo(
            user_role=UserRole.ADMIN,
            user_id=user_id,
            verificated=True,
        ),
        UserInfo(
            user_role=UserRole.ADMIN,
            user_id=user_id,
            verificated=False,
        ),
    )


@fixture(scope="function")
def employee_user_info(user_id) -> tuple[UserInfo, ...]:
    return (
        UserInfo(
            user_role=UserRole.EMPLOYEE,
            user_id=user_id,
            verificated=True,
        ),
        UserInfo(
            user_role=UserRole.EMPLOYEE,
            user_id=user_id,
            verificated=False,
        ),
    )


@fixture(scope="function")
def employer_user_info(user_id) -> tuple[UserInfo, ...]:
    return (
        UserInfo(
            user_role=UserRole.EMPLOYER,
            user_id=user_id,
            verificated=True,
        ),
        UserInfo(
            user_role=UserRole.EMPLOYER,
            user_id=user_id,
            verificated=False,
        ),
    )


@fixture(scope="function")
def unspecified_user_info(user_id) -> tuple[UserInfo, ...]:
    return (
        UserInfo(
            user_role=UserRole.UNSPECIFIED,
            user_id=user_id,
            verificated=True,
        ),
        UserInfo(
            user_role=UserRole.UNSPECIFIED,
            user_id=user_id,
            verificated=False,
        ),
    )


@fixture(scope="function")
def vacancies(user_id) -> tuple[tuple[VacancyORM, UserInfo | None, bool], ...]:
    author_id = uuid4()

    return (
        # published vacancies
        (
            VacancyORM(status=VacancyStatus.PUBLISHED, author_id=author_id),
            UserInfo(user_role=UserRole.ADMIN, user_id=user_id, verificated=True),
            True,
        ),
        (
            VacancyORM(status=VacancyStatus.PUBLISHED, author_id=author_id),
            UserInfo(user_role=UserRole.EMPLOYEE, user_id=user_id, verificated=True),
            True,
        ),
        (
            VacancyORM(status=VacancyStatus.PUBLISHED, author_id=author_id),
            UserInfo(user_role=UserRole.EMPLOYER, user_id=user_id, verificated=True),
            True,
        ),
        (
            VacancyORM(status=VacancyStatus.PUBLISHED, author_id=author_id),
            None,
            True,
        ),
        # moderating vacancies, only author and moderator should see vacancy
        (  # admin
            VacancyORM(status=VacancyStatus.MODERATING, author_id=author_id),
            UserInfo(user_role=UserRole.ADMIN, user_id=user_id, verificated=True),
            True,
        ),
        (  # author
            VacancyORM(status=VacancyStatus.MODERATING, author_id=user_id),
            UserInfo(user_role=UserRole.EMPLOYER, user_id=user_id, verificated=True),
            True,
        ),
        (
            VacancyORM(status=VacancyStatus.MODERATING, author_id=author_id),
            UserInfo(user_role=UserRole.EMPLOYEE, user_id=user_id, verificated=True),
            False,
        ),
        (
            VacancyORM(status=VacancyStatus.MODERATING, author_id=author_id),
            UserInfo(user_role=UserRole.EMPLOYER, user_id=user_id, verificated=True),
            False,
        ),
        (  # no user info
            VacancyORM(status=VacancyStatus.MODERATING, author_id=author_id),
            None,
            False,
        ),
    )
