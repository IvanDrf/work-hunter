from asyncio import wait_for

from grpc.aio import ServicerContext
from src.api.grpc.dependencies import IApplicationService
from src.api.grpc.dto import application_dto
from src.api.grpc.utils import handle_errors, validate_limit_offset

from pkg.applications_api.applications_pb2 import (
    FindVacanciesIDByUserIDRequest,
    FindVacanciesIDByUserIDResponse,
    UpdateApplicationsRequest,
)
from pkg.applications_api.applications_pb2_grpc import ApplicationServiceServicer
from pkg.common.common_pb2 import Response, ResponseStatus, ServiceStatus, Status


class ApplicationHandlers(ApplicationServiceServicer):
    def __init__(self, application_service: IApplicationService, service_timeout: float) -> None:
        super().__init__()

        self.application_service: IApplicationService = application_service
        self.service_timeout: float = service_timeout

    async def stop(self) -> None:
        await self.application_service.stop()

    async def Health(self, request, context) -> ServiceStatus:
        return ServiceStatus(status=Status.AVAILABLE)

    @handle_errors
    async def UpdateApplications(self, request: UpdateApplicationsRequest, context: ServicerContext) -> Response:
        await wait_for(
            self.application_service.update_application(application=application_dto(request)),
            timeout=self.service_timeout,
        )

        return Response(
            status=ResponseStatus.SUCCESS,
            details=f"successfully updated application for user={request.user_info.user_id}, vacancy={request.vacancy_id}",
        )

    @handle_errors
    @validate_limit_offset
    async def FindVacanciesIDByUserID(
        self,
        request: FindVacanciesIDByUserIDRequest,
        context: ServicerContext,
    ) -> FindVacanciesIDByUserIDResponse:
        raise
