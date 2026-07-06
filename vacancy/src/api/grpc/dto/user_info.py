import logging
from uuid import UUID

from pkg.common.common_pb2 import FullUserInfo as PBFullUserInfo
from pkg.common.common_pb2 import UserInfo as PBUserInfo

from src.core.exc import ArgumentError
from src.domain.schemas import UserInfo
from src.domain.types import UserRole


def user_info_dto(user_info: PBUserInfo | PBFullUserInfo) -> UserInfo:
    try:
        return UserInfo(
            user_role=UserRole(user_info.role),
            user_id=UUID(user_info.user_id),
            verificated=user_info.verificated,
        )
    except ValueError as e:
        logging.info(f"Invalid user_id was given: {user_info.user_id}, details={e}")
        raise ArgumentError(f"invalid user_id was given, not uuid: {user_info.user_id}")


def user_info_none_dto(user_info: PBUserInfo | PBFullUserInfo | None) -> UserInfo | None:
    if user_info is None:
        return None

    return user_info_dto(user_info)
