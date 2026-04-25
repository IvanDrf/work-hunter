from grpc import ServicerContext
from pkg.vacancy_api.vacancy_pb2 import CreateVacancyRequest, CreateVacancyResponse, FindVacancyByIDRequest, VacancyInfo
from pkg.vacancy_api.vacancy_pb2_grpc import VacancyServicer

from src.api.dependencies.service import IVacancyService
from src.core.exc.not_found import NotFoundError
from src.utils.handle_errors import handle_errors


class VacancyHandlers(VacancyServicer):
    def __init__(self, vacancy_service: IVacancyService) -> None:
        self.vacancy_service: IVacancyService = vacancy_service
        super().__init__()

    @handle_errors
    async def CreateVacancy(self, request: CreateVacancyRequest, context: ServicerContext) -> CreateVacancyResponse:
        vacancy_id = await self.vacancy_service.create_vacancy(request.vacancy, request.user_info)
        return CreateVacancyResponse(vacancy_id=vacancy_id)

    @handle_errors
    async def FindVacancyByID(self, request: FindVacancyByIDRequest, context: ServicerContext) -> VacancyInfo:
        vacancy = await self.vacancy_service.find_vacancy_by_id(request.vacancy_id, request.user_info)
        if vacancy is None:
            raise NotFoundError(
                f'''can't find vacancy with given {request.vacancy_id=}'''
            )

        return vacancy
