import logging
from asyncio import wait_for

from applications.src.api.grpc.dto.user_info import user_info_dto
from grpc.aio import ServicerContext
from pkg.applications_api.applications_pb2 import (
    FindVacanciesIDByUserIDRequest,
    FindVacanciesIDByUserIDResponse,
    UpdateApplicationsRequest,
)
from pkg.applications_api.applications_pb2_grpc import ApplicationServiceServicer
from pkg.common.common_pb2 import Response, ResponseStatus, ServiceStatus, Status

from src.api.grpc.dependencies import IApplicationService
from src.api.grpc.dto import application_dto
from src.api.grpc.utils import handle_errors, validate_limit_offset


class ApplicationHandlers(ApplicationServiceServicer):
    def __init__(self, application_service: IApplicationService, service_timeout: float) -> None:
        super().__init__()

        self.application_service: IApplicationService = application_service
        self.service_timeout: float = service_timeout

        self.logger = logging.getLogger("ApplicationHandlers")

    async def stop(self) -> None:
        await self.application_service.stop()

    async def Health(self, request, context) -> ServiceStatus:
        return ServiceStatus(status=Status.AVAILABLE)

    @handle_errors
    async def UpdateApplications(self, request: UpdateApplicationsRequest, context: ServicerContext) -> Response:
        self.logger.info(f"UpdateApplications: got request {request=}")
        await wait_for(
            self.application_service.update_application(application=application_dto(request)),
            timeout=self.service_timeout,
        )

        response = Response(
            status=ResponseStatus.SUCCESS,
            details=f"successfully updated application for user={request.user_info.user_id}, vacancy={request.vacancy_id}",
        )
        self.logger.info(f"UpdateApplications: success {request=}{response=}")
        return response

    @handle_errors
    @validate_limit_offset
    async def FindVacanciesIDByUserID(
        self,
        request: FindVacanciesIDByUserIDRequest,
        context: ServicerContext,
    ) -> FindVacanciesIDByUserIDResponse:
        self.logger.info(f"FindVacanciesIDByUserID: got request {request=}")
        vacancies = await wait_for(
            self.application_service.find_vacancies_ids_by_applications(
                user_info_dto(request.user_info),
                limit=request.limit,
                offset=request.offset,
            ),
            timeout=self.service_timeout,
        )

        response = FindVacanciesIDByUserIDResponse(
            user_id=request.user_info.user_id, limit=request.limit, offset=request.offset, vacancies_ids=vacancies
        )
        self.logger.info(f"FindVacanciesIDByUserID: success {request=}{response=}")
        return response
