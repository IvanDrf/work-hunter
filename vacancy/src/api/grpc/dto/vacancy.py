from uuid import UUID

from pkg.vacancy_api.vacancy_pb2 import CreateVacancyRequest, UpdateVacancyRequest, VacancyInfo

from src.domain.schemas import VacancyCreateSchema, VacancyResponseSchema, VacancyUpdateSchema
from src.domain.types import Currency, RemoteType, TimeType


def vacancy_create_dto(request: CreateVacancyRequest) -> VacancyCreateSchema:
    schema = VacancyCreateSchema(
        title=request.title,
        requirements=request.requirements,
        conditions=request.conditions,
        author_id=UUID(request.user_info.user_id),
        author_name=request.user_info.company_name,
        currency=Currency(request.currency),
        remote_type=RemoteType(request.remote_type),
        time_type=TimeType(request.time_type),
        tags=list(request.tags),
    )

    if request.HasField("description"):
        schema.description = request.description

    if request.HasField("city"):
        schema.city = request.city

    if request.HasField("metro"):
        schema.metro = request.metro

    if request.HasField("experience_min"):
        schema.experience_min = request.experience_min

    if request.HasField("experience_max"):
        schema.experience_max = request.experience_max

    if request.HasField("salary_min"):
        schema.salary_min = request.salary_min

    if request.HasField("salary_max"):
        schema.salary_max = request.salary_max

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
        is_city_valid=vacancy.is_city_valid,
        is_metro_valid=vacancy.is_metro_valid,
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


def vacancy_update_dto(request: UpdateVacancyRequest) -> VacancyUpdateSchema:
    schema = VacancyUpdateSchema(vacancy_id=request.vacancy_id)

    if request.HasField("title"):
        schema.title = request.title

    if request.HasField("description"):
        schema.description = request.description

    if request.HasField("requirements"):
        schema.requirements = request.requirements

    if request.HasField("conditions"):
        schema.conditions = request.conditions

    if request.HasField("salary_min"):
        schema.salary_min = request.salary_min if request.salary_min else None

    if request.HasField("salary_max"):
        schema.salary_max = request.salary_max if request.salary_max else None

    if request.HasField("currency"):
        schema.currency = Currency(request.currency)

    if request.HasField("city"):
        schema.city = request.city if request.city else None

    if request.HasField("metro"):
        schema.metro = request.metro if request.metro else None

    if request.HasField("remote_type"):
        schema.remote_type = RemoteType(request.remote_type)

    if request.HasField("time_type"):
        schema.time_type = TimeType(request.time_type)

    if request.HasField("experience_min"):
        schema.experience_min = request.experience_min if request.experience_min else None

    if request.HasField("experience_max"):
        schema.experience_max = request.experience_max if request.experience_max else None

    if request.update_tags:
        schema.tags = list(request.tags)

    return schema
