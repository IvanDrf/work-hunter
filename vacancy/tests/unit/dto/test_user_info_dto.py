from dataclasses import dataclass
from uuid import UUID, uuid4

from pkg.common.common_pb2 import UserInfo as PKGUserInfo
from pkg.common.common_pb2 import UserRole as PKGUserRole
from pkg.vacancy_api.vacancy_pb2 import VacancyStatus as PKGVacancyStatus
from pytest import mark, raises

from src.api.grpc.dto.status import vacancy_status_dto
from src.api.grpc.dto.user_info import user_info_dto, user_info_none_dto
from src.core.exc import ArgumentError
from src.domain.types.enums import UserRole, VacancyStatus


@dataclass(slots=True)
class Request:
    status: PKGVacancyStatus


@mark.parametrize(
    "req",
    [
        Request(PKGVacancyStatus.CLOSED),
        Request(PKGVacancyStatus.PUBLISHED),
        Request(PKGVacancyStatus.DELETED),
        Request(PKGVacancyStatus.MODERATING),
    ],
)
def test_vacancy_status_dto(req: Request) -> None:
    status = vacancy_status_dto(req)

    match req.status:
        case PKGVacancyStatus.CLOSED:
            assert status == VacancyStatus.CLOSED

        case PKGVacancyStatus.PUBLISHED:
            assert status == VacancyStatus.PUBLISHED

        case PKGVacancyStatus.MODERATING:
            assert status == VacancyStatus.MODERATING

        case PKGVacancyStatus.DELETED:
            assert status == VacancyStatus.DELETED


@mark.parametrize(
    ["role", "user_id", "verificated", "is_user_id_valid"],
    [
        (PKGUserRole.ADMIN, str(uuid4()), True, True),
        (PKGUserRole.ADMIN, str(uuid4()), False, True),
        (PKGUserRole.EMPLOYEE, str(uuid4()), True, True),
        (PKGUserRole.EMPLOYEE, str(uuid4()), False, True),
        (PKGUserRole.EMPLOYER, str(uuid4()), True, True),
        (PKGUserRole.EMPLOYER, str(uuid4()), False, True),
        # invalid user_id - not UUID
        (PKGUserRole.ADMIN, "12", True, False),
        (PKGUserRole.ADMIN, "user_id", False, False),
        (PKGUserRole.EMPLOYEE, "rlkweef", True, False),
        (PKGUserRole.EMPLOYEE, "id:12345", False, False),
        (PKGUserRole.EMPLOYER, "rk;gw", True, False),
        (PKGUserRole.EMPLOYER, "-20192371", False, False),
    ],
)
def test_user_info_dto(role: PKGUserRole, user_id: str, verificated: bool, is_user_id_valid: bool) -> None:
    roles = {
        PKGUserRole.ADMIN: UserRole.ADMIN,
        PKGUserRole.EMPLOYEE: UserRole.EMPLOYEE,
        PKGUserRole.EMPLOYER: UserRole.EMPLOYER,
    }

    pkg_user_info = PKGUserInfo(role=role, user_id=user_id, verificated=verificated)
    if not is_user_id_valid:
        error_message = f"invalid user_id was given: {user_id}"

        with raises(ArgumentError, match=error_message):
            user_info_dto(pkg_user_info)
    else:
        user_info = user_info_dto(pkg_user_info)
        assert user_info.user_id == UUID(user_id)
        assert user_info.verificated == verificated
        assert user_info.user_role == roles[role]


@mark.parametrize(
    ["role", "user_id", "verificated", "is_user_id_valid"],
    [
        (PKGUserRole.ADMIN, str(uuid4()), True, True),
        (PKGUserRole.ADMIN, str(uuid4()), False, True),
        (PKGUserRole.EMPLOYEE, str(uuid4()), True, True),
        (PKGUserRole.EMPLOYEE, str(uuid4()), False, True),
        (PKGUserRole.EMPLOYER, str(uuid4()), True, True),
        (PKGUserRole.EMPLOYER, str(uuid4()), False, True),
        # invalid user_id - not UUID
        (PKGUserRole.ADMIN, "12", True, False),
        (PKGUserRole.ADMIN, "user_id", False, False),
        (PKGUserRole.EMPLOYEE, "rlkweef", True, False),
        (PKGUserRole.EMPLOYEE, "id:12345", False, False),
        (PKGUserRole.EMPLOYER, "rk;gw", True, False),
        (PKGUserRole.EMPLOYER, "-20192371", False, False),
    ],
)
def test_user_info_none_dto(role: PKGUserRole, user_id: str, verificated: bool, is_user_id_valid: bool) -> None:
    roles = {
        PKGUserRole.ADMIN: UserRole.ADMIN,
        PKGUserRole.EMPLOYEE: UserRole.EMPLOYEE,
        PKGUserRole.EMPLOYER: UserRole.EMPLOYER,
    }

    pkg_user_info = PKGUserInfo(role=role, user_id=user_id, verificated=verificated)
    if not is_user_id_valid:
        error_message = f"invalid user_id was given: {user_id}"

        with raises(ArgumentError, match=error_message):
            user_info_none_dto(pkg_user_info)
    else:
        user_info = user_info_none_dto(pkg_user_info)
        assert user_info is not None
        assert user_info.user_id == UUID(user_id)
        assert user_info.verificated == verificated
        assert user_info.user_role == roles[role]

    assert user_info_none_dto(None) is None
