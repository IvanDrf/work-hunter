from uuid import UUID

from pkg.common.common_pb2 import UserInfo


def is_user_id_valid(user_info: UserInfo) -> bool:
    try:
        UUID(user_info.user_id)
        return True
    except ValueError:
        return False


def get_user_info(request) -> UserInfo | None:
    return request.user_info if request.HasField('user_info') else None
