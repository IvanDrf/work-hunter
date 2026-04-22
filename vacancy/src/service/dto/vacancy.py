from datetime import datetime, timezone

from src.domain.models.vacancy import VacancyORM
from pkg.vacancy_api.vacancy_pb2 import VacancyInfo
from pkg.common.common_pb2 import UserInfo


def create_vacancy_dto(vacancy: VacancyInfo, user_info: UserInfo) -> VacancyORM:
    return VacancyORM(
        author_id=user_info.user_id,

        title=vacancy.title,
        description=vacancy.description,

        requirements=vacancy.requirements,
        conditions=vacancy.conditions,

        salary_min=vacancy.salary_min,
        salary_max=vacancy.salary_max,
        currenct=vacancy.currency,

        city=None if vacancy.city == '' else vacancy.city,
        metro=None if vacancy.metro == '' else vacancy.metro,

        remote_type=vacancy.remote_type,
        time_type=vacancy.time_type,

        experience_min=None if vacancy.experience_min == 0 else vacancy.experience_min,
        experience_max=None if vacancy.experience_max == 0 else vacancy.experience_max,

        created_at=datetime.now(timezone.utc),
        updated_at=None,
        published_at=None,
        closed_at=None,

        status=vacancy.status,
        moderated_at=None,
        moderator_comments=None,

        views=0,
        applications_count=0,
    )
