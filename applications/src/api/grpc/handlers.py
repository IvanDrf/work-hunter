from grpc.aio import ServicerContext

from pkg.applications_api.applications_pb2 import UpdateApplicationsRequest
from pkg.applications_api.applications_pb2_grpc import ApplicationServiceServicer
from pkg.common.common_pb2 import Response, ServiceStatus, Status


class Handlers(ApplicationServiceServicer):
    async def Health(self, request, context) -> ServiceStatus:
        return ServiceStatus(status=Status.AVAILABLE)

    async def UpdateApplications(self, request: UpdateApplicationsRequest, context: ServicerContext) -> Response:
        
