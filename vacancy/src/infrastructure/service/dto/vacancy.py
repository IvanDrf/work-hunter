from datetime import datetime

from src.domain.models.tag import TagORM
from src.domain.models.vacancy import VacancyORM, VacancyStatus
from src.domain.schemas import UserInfo, VacancyCreateSchema, VacancyResponseSchema


def create_vacancy_dto(
    vacancy: VacancyCreateSchema, user_info: UserInfo, created_at: datetime, status: VacancyStatus
) -> VacancyORM:
    return VacancyORM(
        author_id=user_info.user_id,
        author_name=vacancy.author_name,
        title=vacancy.title,
        description=vacancy.description,
        requirements=vacancy.requirements,
        conditions=vacancy.conditions,
        salary_min=vacancy.salary_min,
        salary_max=vacancy.salary_max,
        currency=vacancy.currency,
        city=None if vacancy.city == "" else vacancy.city,
        metro=None if vacancy.metro == "" else vacancy.metro,
        remote_type=vacancy.remote_type,
        time_type=vacancy.time_type,
        experience_min=None if vacancy.experience_min == 0 else vacancy.experience_min,
        experience_max=None if vacancy.experience_max == 0 else vacancy.experience_max,
        created_at=created_at,
        updated_at=None,
        published_at=None,
        closed_at=None,
        status=status,
        moderated_at=None,
        moderator_comments=None,
        views=0,
        applications_count=0,
        tags=[TagORM(tag=t) for t in vacancy.tags] if vacancy.tags is not None else [],
    )


def vacancy_orm_to_response_dto(vacancy: VacancyORM) -> VacancyResponseSchema:
    return VacancyResponseSchema(
        author_name=vacancy.author_name,
        vacancy_id=vacancy.vacancy_id,
        title=vacancy.title,
        description=vacancy.description,
        requirements=vacancy.requirements,
        conditions=vacancy.conditions,
        salary_min=vacancy.salary_min,
        salary_max=vacancy.salary_max,
        currency=vacancy.currency,
        city=vacancy.city,
        metro=vacancy.metro,
        remote_type=vacancy.remote_type,
        time_type=vacancy.time_type,
        experience_min=vacancy.experience_min,
        experience_max=vacancy.experience_max,
        created_at=vacancy.created_at,
        updated_at=vacancy.updated_at,
        published_at=vacancy.published_at,
        closed_at=vacancy.closed_at,
        status=vacancy.status,
        moderated_at=vacancy.moderated_at,
        moderator_comments=vacancy.moderator_comments,
        views=vacancy.views,
        applications_count=vacancy.applications_count,
        tags=[tag.tag for tag in vacancy.tags],
    )
