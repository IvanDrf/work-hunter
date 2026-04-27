import datetime

from google.protobuf import timestamp_pb2 as _timestamp_pb2
from ..common import common_pb2 as _common_pb2
from google.protobuf.internal import containers as _containers
from google.protobuf.internal import enum_type_wrapper as _enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor


class RemoteType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    OFFICE: _ClassVar[RemoteType]
    REMOTE: _ClassVar[RemoteType]
    HYBRID: _ClassVar[RemoteType]
    ANY: _ClassVar[RemoteType]


class TimeType(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    FULL: _ClassVar[TimeType]
    PART: _ClassVar[TimeType]


class Currency(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    RUB: _ClassVar[Currency]
    USD: _ClassVar[Currency]
    EUR: _ClassVar[Currency]


class VacancyStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    MODERATING: _ClassVar[VacancyStatus]
    PUBLISHED: _ClassVar[VacancyStatus]
    UPDATED: _ClassVar[VacancyStatus]
    CLOSED: _ClassVar[VacancyStatus]
    DELETED: _ClassVar[VacancyStatus]


class ResponseStatus(int, metaclass=_enum_type_wrapper.EnumTypeWrapper):
    __slots__ = ()
    SUCCESS: _ClassVar[ResponseStatus]
    FAILED: _ClassVar[ResponseStatus]
    NOT_FOUND: _ClassVar[ResponseStatus]
    FORBIDDEN: _ClassVar[ResponseStatus]


OFFICE: RemoteType
REMOTE: RemoteType
HYBRID: RemoteType
ANY: RemoteType
FULL: TimeType
PART: TimeType
RUB: Currency
USD: Currency
EUR: Currency
MODERATING: VacancyStatus
PUBLISHED: VacancyStatus
UPDATED: VacancyStatus
CLOSED: VacancyStatus
DELETED: VacancyStatus
SUCCESS: ResponseStatus
FAILED: ResponseStatus
NOT_FOUND: ResponseStatus
FORBIDDEN: ResponseStatus


class Response(_message.Message):
    __slots__ = ("message", "status")
    MESSAGE_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    message: str
    status: ResponseStatus
    def __init__(self, message: _Optional[str] = ...,
                 status: _Optional[_Union[ResponseStatus, str]] = ...) -> None: ...


class VacancyInfo(_message.Message):
    __slots__ = ("vacancy_id", "title", "description", "requirements", "conditions", "salary_min", "salary_max", "currency", "city", "metro", "remote_type", "time_type", "experience_min",
                 "experience_max", "created_at", "updated_at", "published_at", "closed_at", "status", "moderated_time", "moderator_comments", "views", "applications_count", "tags")
    VACANCY_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    REQUIREMENTS_FIELD_NUMBER: _ClassVar[int]
    CONDITIONS_FIELD_NUMBER: _ClassVar[int]
    SALARY_MIN_FIELD_NUMBER: _ClassVar[int]
    SALARY_MAX_FIELD_NUMBER: _ClassVar[int]
    CURRENCY_FIELD_NUMBER: _ClassVar[int]
    CITY_FIELD_NUMBER: _ClassVar[int]
    METRO_FIELD_NUMBER: _ClassVar[int]
    REMOTE_TYPE_FIELD_NUMBER: _ClassVar[int]
    TIME_TYPE_FIELD_NUMBER: _ClassVar[int]
    EXPERIENCE_MIN_FIELD_NUMBER: _ClassVar[int]
    EXPERIENCE_MAX_FIELD_NUMBER: _ClassVar[int]
    CREATED_AT_FIELD_NUMBER: _ClassVar[int]
    UPDATED_AT_FIELD_NUMBER: _ClassVar[int]
    PUBLISHED_AT_FIELD_NUMBER: _ClassVar[int]
    CLOSED_AT_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    MODERATED_TIME_FIELD_NUMBER: _ClassVar[int]
    MODERATOR_COMMENTS_FIELD_NUMBER: _ClassVar[int]
    VIEWS_FIELD_NUMBER: _ClassVar[int]
    APPLICATIONS_COUNT_FIELD_NUMBER: _ClassVar[int]
    TAGS_FIELD_NUMBER: _ClassVar[int]
    vacancy_id: int
    title: str
    description: str
    requirements: str
    conditions: str
    salary_min: int
    salary_max: int
    currency: Currency
    city: str
    metro: str
    remote_type: RemoteType
    time_type: TimeType
    experience_min: int
    experience_max: int
    created_at: _timestamp_pb2.Timestamp
    updated_at: _timestamp_pb2.Timestamp
    published_at: _timestamp_pb2.Timestamp
    closed_at: _timestamp_pb2.Timestamp
    status: VacancyStatus
    moderated_time: _timestamp_pb2.Timestamp
    moderator_comments: str
    views: int
    applications_count: int
    tags: _containers.RepeatedScalarFieldContainer[str]
    def __init__(self, vacancy_id: _Optional[int] = ..., title: _Optional[str] = ..., description: _Optional[str] = ..., requirements: _Optional[str] = ..., conditions: _Optional[str] = ..., salary_min: _Optional[int] = ..., salary_max: _Optional[int] = ..., currency: _Optional[_Union[Currency, str]] = ..., city: _Optional[str] = ..., metro: _Optional[str] = ..., remote_type: _Optional[_Union[RemoteType, str]] = ..., time_type: _Optional[_Union[TimeType, str]] = ..., experience_min: _Optional[int] = ..., experience_max: _Optional[int] = ..., created_at: _Optional[_Union[datetime.datetime, _timestamp_pb2.Timestamp,
                 _Mapping]] = ..., updated_at: _Optional[_Union[datetime.datetime, _timestamp_pb2.Timestamp, _Mapping]] = ..., published_at: _Optional[_Union[datetime.datetime, _timestamp_pb2.Timestamp, _Mapping]] = ..., closed_at: _Optional[_Union[datetime.datetime, _timestamp_pb2.Timestamp, _Mapping]] = ..., status: _Optional[_Union[VacancyStatus, str]] = ..., moderated_time: _Optional[_Union[datetime.datetime, _timestamp_pb2.Timestamp, _Mapping]] = ..., moderator_comments: _Optional[str] = ..., views: _Optional[int] = ..., applications_count: _Optional[int] = ..., tags: _Optional[_Iterable[str]] = ...) -> None: ...


class CreateVacancyRequest(_message.Message):
    __slots__ = ("vacancy", "user_info")
    VACANCY_FIELD_NUMBER: _ClassVar[int]
    USER_INFO_FIELD_NUMBER: _ClassVar[int]
    vacancy: VacancyInfo
    user_info: _common_pb2.UserInfo
    def __init__(self, vacancy: _Optional[_Union[VacancyInfo, _Mapping]] = ...,
                 user_info: _Optional[_Union[_common_pb2.UserInfo, _Mapping]] = ...) -> None: ...


class CreateVacancyResponse(_message.Message):
    __slots__ = ("vacancy_id",)
    VACANCY_ID_FIELD_NUMBER: _ClassVar[int]
    vacancy_id: int
    def __init__(self, vacancy_id: _Optional[int] = ...) -> None: ...


class UpdateVacancyRequest(_message.Message):
    __slots__ = ("vacancy_id", "title", "description", "requirements", "conditions", "salary_min", "salary_max", "currency", "city", "metro", "remote_type", "time_type", "experience_min",
                 "experience_max", "created_at", "updated_at", "published_at", "closed_at", "status", "moderated_time", "moderator_comments", "views", "applications_count", "tags", "user_info")
    VACANCY_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    REQUIREMENTS_FIELD_NUMBER: _ClassVar[int]
    CONDITIONS_FIELD_NUMBER: _ClassVar[int]
    SALARY_MIN_FIELD_NUMBER: _ClassVar[int]
    SALARY_MAX_FIELD_NUMBER: _ClassVar[int]
    CURRENCY_FIELD_NUMBER: _ClassVar[int]
    CITY_FIELD_NUMBER: _ClassVar[int]
    METRO_FIELD_NUMBER: _ClassVar[int]
    REMOTE_TYPE_FIELD_NUMBER: _ClassVar[int]
    TIME_TYPE_FIELD_NUMBER: _ClassVar[int]
    EXPERIENCE_MIN_FIELD_NUMBER: _ClassVar[int]
    EXPERIENCE_MAX_FIELD_NUMBER: _ClassVar[int]
    CREATED_AT_FIELD_NUMBER: _ClassVar[int]
    UPDATED_AT_FIELD_NUMBER: _ClassVar[int]
    PUBLISHED_AT_FIELD_NUMBER: _ClassVar[int]
    CLOSED_AT_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    MODERATED_TIME_FIELD_NUMBER: _ClassVar[int]
    MODERATOR_COMMENTS_FIELD_NUMBER: _ClassVar[int]
    VIEWS_FIELD_NUMBER: _ClassVar[int]
    APPLICATIONS_COUNT_FIELD_NUMBER: _ClassVar[int]
    TAGS_FIELD_NUMBER: _ClassVar[int]
    USER_INFO_FIELD_NUMBER: _ClassVar[int]
    vacancy_id: int
    title: str
    description: str
    requirements: str
    conditions: str
    salary_min: int
    salary_max: int
    currency: Currency
    city: str
    metro: str
    remote_type: RemoteType
    time_type: TimeType
    experience_min: int
    experience_max: int
    created_at: _timestamp_pb2.Timestamp
    updated_at: _timestamp_pb2.Timestamp
    published_at: _timestamp_pb2.Timestamp
    closed_at: _timestamp_pb2.Timestamp
    status: VacancyStatus
    moderated_time: _timestamp_pb2.Timestamp
    moderator_comments: str
    views: int
    applications_count: int
    tags: _containers.RepeatedScalarFieldContainer[str]
    user_info: _common_pb2.UserInfo
    def __init__(self, vacancy_id: _Optional[int] = ..., title: _Optional[str] = ..., description: _Optional[str] = ..., requirements: _Optional[str] = ..., conditions: _Optional[str] = ..., salary_min: _Optional[int] = ..., salary_max: _Optional[int] = ..., currency: _Optional[_Union[Currency, str]] = ..., city: _Optional[str] = ..., metro: _Optional[str] = ..., remote_type: _Optional[_Union[RemoteType, str]] = ..., time_type: _Optional[_Union[TimeType, str]] = ..., experience_min: _Optional[int] = ..., experience_max: _Optional[int] = ..., created_at: _Optional[_Union[datetime.datetime, _timestamp_pb2.Timestamp, _Mapping]] = ...,
                 updated_at: _Optional[_Union[datetime.datetime, _timestamp_pb2.Timestamp, _Mapping]] = ..., published_at: _Optional[_Union[datetime.datetime, _timestamp_pb2.Timestamp, _Mapping]] = ..., closed_at: _Optional[_Union[datetime.datetime, _timestamp_pb2.Timestamp, _Mapping]] = ..., status: _Optional[_Union[VacancyStatus, str]] = ..., moderated_time: _Optional[_Union[datetime.datetime, _timestamp_pb2.Timestamp, _Mapping]] = ..., moderator_comments: _Optional[str] = ..., views: _Optional[int] = ..., applications_count: _Optional[int] = ..., tags: _Optional[_Iterable[str]] = ..., user_info: _Optional[_Union[_common_pb2.UserInfo, _Mapping]] = ...) -> None: ...


class DeleteVacancyRequest(_message.Message):
    __slots__ = ("vacancy_id", "user_info")
    VACANCY_ID_FIELD_NUMBER: _ClassVar[int]
    USER_INFO_FIELD_NUMBER: _ClassVar[int]
    vacancy_id: int
    user_info: _common_pb2.UserInfo
    def __init__(self, vacancy_id: _Optional[int] = ...,
                 user_info: _Optional[_Union[_common_pb2.UserInfo, _Mapping]] = ...) -> None: ...


class FindVacancyByIDRequest(_message.Message):
    __slots__ = ("vacancy_id", "user_info")
    VACANCY_ID_FIELD_NUMBER: _ClassVar[int]
    USER_INFO_FIELD_NUMBER: _ClassVar[int]
    vacancy_id: int
    user_info: _common_pb2.UserInfo
    def __init__(self, vacancy_id: _Optional[int] = ...,
                 user_info: _Optional[_Union[_common_pb2.UserInfo, _Mapping]] = ...) -> None: ...


class FindVacancyByTagsRequest(_message.Message):
    __slots__ = ("tags", "limit", "offset", "user_info")
    TAGS_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    OFFSET_FIELD_NUMBER: _ClassVar[int]
    USER_INFO_FIELD_NUMBER: _ClassVar[int]
    tags: _containers.RepeatedScalarFieldContainer[str]
    limit: int
    offset: int
    user_info: _common_pb2.UserInfo
    def __init__(self, tags: _Optional[_Iterable[str]] = ..., limit: _Optional[int] = ..., offset: _Optional[int]
                 = ..., user_info: _Optional[_Union[_common_pb2.UserInfo, _Mapping]] = ...) -> None: ...


class Vacancies(_message.Message):
    __slots__ = ("vacancies", "limit", "offset")
    VACANCIES_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    OFFSET_FIELD_NUMBER: _ClassVar[int]
    vacancies: _containers.RepeatedCompositeFieldContainer[VacancyInfo]
    limit: int
    offset: int
    def __init__(self, vacancies: _Optional[_Iterable[_Union[VacancyInfo, _Mapping]]]
                 = ..., limit: _Optional[int] = ..., offset: _Optional[int] = ...) -> None: ...


class FindVacanciesByAuthorRequest(_message.Message):
    __slots__ = ("author", "user_info")
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    USER_INFO_FIELD_NUMBER: _ClassVar[int]
    author: str
    user_info: _common_pb2.UserInfo
    def __init__(self, author: _Optional[str] = ...,
                 user_info: _Optional[_Union[_common_pb2.UserInfo, _Mapping]] = ...) -> None: ...


class SetVacancyStatusRequest(_message.Message):
    __slots__ = ("vacancy_id", "status", "user_info")
    VACANCY_ID_FIELD_NUMBER: _ClassVar[int]
    STATUS_FIELD_NUMBER: _ClassVar[int]
    USER_INFO_FIELD_NUMBER: _ClassVar[int]
    vacancy_id: int
    status: VacancyStatus
    user_info: _common_pb2.UserInfo

    def __init__(self, vacancy_id: _Optional[int] = ..., status: _Optional[_Union[VacancyStatus, str]]
                 = ..., user_info: _Optional[_Union[_common_pb2.UserInfo, _Mapping]] = ...) -> None: ...
