from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class UserRole(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNSPECIFIED: _ClassVar[UserRole]
    ADMIN: _ClassVar[UserRole]
    EMPLOYEE: _ClassVar[UserRole]
    EMPLOYER: _ClassVar[UserRole]

class ResponseStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    SUCCESS: _ClassVar[ResponseStatus]
    FAILED: _ClassVar[ResponseStatus]
    NOT_FOUND: _ClassVar[ResponseStatus]
    FORBIDDEN: _ClassVar[ResponseStatus]

class Status(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    UNAVAILABLE: _ClassVar[Status]
    AVAILABLE: _ClassVar[Status]
UNSPECIFIED: UserRole
ADMIN: UserRole
EMPLOYEE: UserRole
EMPLOYER: UserRole
SUCCESS: ResponseStatus
FAILED: ResponseStatus
NOT_FOUND: ResponseStatus
FORBIDDEN: ResponseStatus
UNAVAILABLE: Status
AVAILABLE: Status

class UserInfo(_message.Message):
    __slots__ = ("role", "user_id", "verificated")
    ROLE_FIELD_NUMBER: _ClassVar[int]
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    VERIFICATED_FIELD_NUMBER: _ClassVar[int]
    role: UserRole
    user_id: str
    verificated: bool
    def __init__(self, role: _Optional[_Union[UserRole, str]] = ..., user_id: _Optional[str] = ..., verificated: _Optional[bool] = ...) -> None: ...

class FullUserInfo(_message.Message):
    __slots__ = ("role", "user_id", "company_name", "verificated")
    ROLE_FIELD_NUMBER: _ClassVar[int]
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    COMPANY_NAME_FIELD_NUMBER: _ClassVar[int]
    VERIFICATED_FIELD_NUMBER: _ClassVar[int]
    role: UserRole
    user_id: str
    company_name: str
    verificated: bool
    def __init__(self, role: _Optional[_Union[UserRole, str]] = ..., user_id: _Optional[str] = ..., company_name: _Optional[str] = ..., verificated: _Optional[bool] = ...) -> None: ...

class Response(_message.Message):
    __slots__ = ("status", "details")
    STATUS_FIELD_NUMBER: _ClassVar[int]
    DETAILS_FIELD_NUMBER: _ClassVar[int]
    status: ResponseStatus
    details: str
    def __init__(self, status: _Optional[_Union[ResponseStatus, str]] = ..., details: _Optional[str] = ...) -> None: ...

class Empty(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class ServiceStatus(_message.Message):
    __slots__ = ("status",)
    STATUS_FIELD_NUMBER: _ClassVar[int]
    status: Status
    def __init__(self, status: _Optional[_Union[Status, str]] = ...) -> None: ...
