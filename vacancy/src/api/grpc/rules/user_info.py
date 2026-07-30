from pkg.common.common_pb2 import UserInfo


def get_user_info(request) -> UserInfo | None:
    return request.user_info if request.HasField("user_info") else None
