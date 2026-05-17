from uuid import UUID

from pkg.common.common_pb2 import FullUserInfo
from pkg.vacancy_api.vacancy_pb2 import UpdateVacancyRequest, VacancyInfo

from src.domain.schemas import VacancyCreateSchema, VacancyResponseSchema, VacancyUpdateSchema
from src.domain.types.enums import Currency, RemoteType, TimeType
from src.domain.types.types import UNSET_VALUE, Money, Year


def vacancy_create_dto(vacancy: VacancyInfo, user_info: FullUserInfo) -> VacancyCreateSchema:
    schema = VacancyCreateSchema(
        title=vacancy.title,
        requirements=vacancy.requirements,
        conditions=vacancy.requirements,
        author_id=UUID(user_info.user_id),
        author_name=user_info.username,
        salary_min=vacancy.salary_min,
        salary_max=vacancy.salary_max,
        currency=Currency(vacancy.currency),
        remote_type=RemoteType(vacancy.remote_type),
        time_type=TimeType(vacancy.time_type),
        tags=list(vacancy.tags),
    )

    if vacancy.description:
        schema.description = vacancy.description

    if vacancy.city:
        schema.city = vacancy.city

    if vacancy.metro:
        schema.metro = vacancy.metro

    if vacancy.experience_min:
        schema.experience_min = vacancy.experience_min

    if vacancy.experience_max:
        schema.experience_max = vacancy.experience_max

    return schema


def vacancy_response_dto(vacancy: VacancyResponseSchema) -> VacancyInfo:
    return VacancyInfo(
        vacancy_id=vacancy.vacancy_id,
        title=vacancy.title,
        description=vacancy.description,
        requirements=vacancy.requirements,
        conditions=vacancy.conditions,
        salary_min=vacancy.salary_min,
        salary_max=vacancy.salary_max,
        currency=vacancy.currency.name,
        experience_min=vacancy.experience_min,
        experience_max=vacancy.experience_max,
        created_at=vacancy.created_at,
        status=vacancy.status.name,
        remote_type=vacancy.remote_type.name,
        time_type=vacancy.time_type.name,
        city=vacancy.city,
        metro=vacancy.metro,
        views=vacancy.views,
        applications_count=vacancy.applications_count,
        tags=vacancy.tags,
        author_name=vacancy.author_name,
        moderated_time=vacancy.moderated_at,
        moderator_comments=vacancy.moderator_comments,
        updated_at=vacancy.updated_at,
        published_at=vacancy.published_at,
        closed_at=vacancy.closed_at,
    )


def vacancy_update_dto(vacancy: UpdateVacancyRequest) -> VacancyUpdateSchema:
    return VacancyUpdateSchema(
        vacancy_id=vacancy.vacancy_id,
        title=vacancy.title if vacancy.title else UNSET_VALUE,
        description=vacancy.description if vacancy.description else UNSET_VALUE,
        requirements=vacancy.requirements if vacancy.requirements else UNSET_VALUE,
        conditions=vacancy.conditions if vacancy.conditions else UNSET_VALUE,
        salary_min=Money(vacancy.salary_min) if vacancy.salary_min else UNSET_VALUE,
        salary_max=Money(vacancy.salary_max) if vacancy.salary_max else UNSET_VALUE,
        currency=Currency(vacancy.currency) if vacancy.currency else UNSET_VALUE,
        city=vacancy.city if vacancy.city else UNSET_VALUE,
        metro=vacancy.metro if vacancy.metro else UNSET_VALUE,
        remote_type=RemoteType(vacancy.remote_type) if vacancy.remote_type else UNSET_VALUE,
        time_type=TimeType(vacancy.time_type) if vacancy.time_type else UNSET_VALUE,
        experience_min=Year(vacancy.experience_min) if vacancy.experience_min else UNSET_VALUE,
        experience_max=Year(vacancy.experience_max) if vacancy.experience_max else UNSET_VALUE,
        tags=list(vacancy.tags) if vacancy.tags else UNSET_VALUE,
    )
