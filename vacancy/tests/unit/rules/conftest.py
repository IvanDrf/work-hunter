from uuid import UUID

from pytest import fixture

from src.domain.schemas.user import UserInfo
from src.domain.types.enums import UserRole


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
