from datetime import datetime, timezone

from pkg.common.common_pb2 import UserInfo
from pkg.vacancy_api.vacancy_pb2 import Currency as PBCurrency
from pkg.vacancy_api.vacancy_pb2 import RemoteType as PBRemoteType
from pkg.vacancy_api.vacancy_pb2 import TimeType as PBTimeType
from pkg.vacancy_api.vacancy_pb2 import VacancyInfo
from pkg.vacancy_api.vacancy_pb2 import VacancyStatus as PBVacancyStatus

from src.domain.models.vacancy import VacancyORM, VacancyStatus


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

        status=VacancyStatus.MODERATING,
        moderated_at=None,
        moderator_comments=None,

        views=0,
        applications_count=0,
    )


def vacancy_info_dto(vacancy: VacancyORM) -> VacancyInfo:
    return VacancyInfo(
        vacancy_id=vacancy.vacancy_id,
        title=vacancy.title,
        description=vacancy.description,

        requirements=vacancy.requirements,
        conditions=vacancy.conditions,

        salary_min=vacancy.salary_min,
        salary_max=vacancy.salary_max,
        currency=PBCurrency(vacancy.currency.value),

        city=vacancy.city,
        metro=vacancy.metro,

        remote_type=PBRemoteType(vacancy.remote_type.value),
        time_type=PBTimeType(vacancy.time_type.value),

        experience_min=vacancy.experience_min,
        experience_max=vacancy.experience_max,

        created_at=vacancy.created_at,
        updated_at=vacancy.updated_at,
        published_at=vacancy.published_at,
        closed_at=vacancy.closed_at,

        status=PBVacancyStatus(vacancy.status.value),
        moderated_time=vacancy.moderated_at,
        moderator_comments=vacancy.moderator_comments,

        views=vacancy.views,
        applications_count=vacancy.applications_count,
    )
