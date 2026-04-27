from typing import ClassVar as _ClassVar
from typing import Optional as _Optional
from typing import Union as _Union

from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper

DESCRIPTOR: _descriptor.FileDescriptor

class UserRole(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    ADMIN: _ClassVar[UserRole]
    EMPLOYEE: _ClassVar[UserRole]
    EMPLOYER: _ClassVar[UserRole]
ADMIN: UserRole
EMPLOYEE: UserRole
EMPLOYER: UserRole

class UserInfo(_message.Message):
    __slots__ = ("role", "user_id", "verificated")
    ROLE_FIELD_NUMBER: _ClassVar[int]
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    VERIFICATED_FIELD_NUMBER: _ClassVar[int]
    role: UserRole
    user_id: str
    verificated: bool
    def __init__(self, role: _Optional[_Union[UserRole, str]] = ..., user_id: _Optional[str] = ..., verificated: bool = ...) -> None: ...
