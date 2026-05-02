from pkg.common.common_pb2 import UserInfo, UserRole


def is_user_admin(user_info: UserInfo) -> bool:
    return user_info.role == UserRole.ADMIN


def is_user_employer(user_info: UserInfo) -> bool:
    return user_info.role == UserRole.EMPLOYER
