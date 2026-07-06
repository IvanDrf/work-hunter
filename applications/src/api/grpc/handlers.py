from pkg.applications_api.applications_pb2_grpc import ApplicationServiceServicer
from pkg.common.common_pb2 import ServiceStatus, Status


class Handlers(ApplicationServiceServicer):
    async def Health(self, request, context) -> ServiceStatus:
        return ServiceStatus(status=Status.AVAILABLE)
