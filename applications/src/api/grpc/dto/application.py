from src.api.grpc.dto.user_info import user_info_dto
from src.domain.schemas import ApplicationSchema

from pkg.applications_api.applications_pb2 import UpdateApplicationsRequest


def application_dto(application: UpdateApplicationsRequest) -> ApplicationSchema:
    return ApplicationSchema(vacancy_id=application.vacancy_id, user_info=user_info_dto(application.user_info))
