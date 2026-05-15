from uuid import UUID

from pkg.common.common_pb2 import FullUserInfo
from pkg.vacancy_api.vacancy_pb2 import VacancyInfo

from src.domain.schemas import VacancyCreateSchema, VacancyResponseSchema
from src.domain.types.enums import Currency, RemoteType, TimeType


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
