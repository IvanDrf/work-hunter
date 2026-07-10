from src.domain.models import ApplicationORM
from src.domain.schemas import ApplicationSchema


def application_dto(application: ApplicationSchema) -> ApplicationORM:
    return ApplicationORM(
        vacancy_id=application.vacancy_id,
        user_id=application.user_info.user_id,
    )
