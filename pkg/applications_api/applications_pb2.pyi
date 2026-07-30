from ..common import common_pb2 as _common_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class UpdateApplicationsRequest(_message.Message):
    __slots__ = ("vacancy_id", "user_info")
    VACANCY_ID_FIELD_NUMBER: _ClassVar[int]
    USER_INFO_FIELD_NUMBER: _ClassVar[int]
    vacancy_id: int
    user_info: _common_pb2.UserInfo
    def __init__(
        self, vacancy_id: _Optional[int] = ..., user_info: _Optional[_Union[_common_pb2.UserInfo, _Mapping]] = ...
    ) -> None: ...

class FindVacanciesIDByUserIDRequest(_message.Message):
    __slots__ = ("user_id", "limit", "offset", "user_info")
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    OFFSET_FIELD_NUMBER: _ClassVar[int]
    USER_INFO_FIELD_NUMBER: _ClassVar[int]
    user_id: str
    limit: int
    offset: int
    user_info: _common_pb2.UserInfo
    def __init__(
        self,
        user_id: _Optional[str] = ...,
        limit: _Optional[int] = ...,
        offset: _Optional[int] = ...,
        user_info: _Optional[_Union[_common_pb2.UserInfo, _Mapping]] = ...,
    ) -> None: ...

class FindVacanciesIDByUserIDResponse(_message.Message):
    __slots__ = ("user_id", "limit", "offset", "vacancies_ids")
    USER_ID_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    OFFSET_FIELD_NUMBER: _ClassVar[int]
    VACANCIES_IDS_FIELD_NUMBER: _ClassVar[int]
    user_id: str
    limit: int
    offset: int
    vacancies_ids: _containers.RepeatedScalarFieldContainer[int]
    def __init__(
        self,
        user_id: _Optional[str] = ...,
        limit: _Optional[int] = ...,
        offset: _Optional[int] = ...,
        vacancies_ids: _Optional[_Iterable[int]] = ...,
    ) -> None: ...
