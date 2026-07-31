import logging
from uuid import UUID

from pkg.common.common_pb2 import UserInfo as PKGUserInfo

from src.core.exc import ArgumentError
from src.domain.schemas import UserInfo, UserRole


def user_info_dto(user_info: PKGUserInfo) -> UserInfo:
    try:
        return UserInfo(
            user_role=UserRole(user_info.role),
            user_id=UUID(user_info.user_id),
            verificated=user_info.verificated,
        )

    except ValueError as e:
        logging.getLogger("user_info_dto").info(f"Invalid user_id was given: {user_info.user_id}, details={e}")
        raise ArgumentError(f"invalid user_id was given, not uuid: {user_info.user_id}")
