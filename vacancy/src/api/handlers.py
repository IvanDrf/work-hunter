import logging

from grpc import ServicerContext, StatusCode
from pkg.vacancy_api.vacancy_pb2 import CreateVacancyRequest, CreateVacancyResponse
from pkg.vacancy_api.vacancy_pb2_grpc import VacancyServicer

from src.api.dependencies import IVacancyService
from src.core.exc.internal import InternalError
from src.core.exc.invalid_argument import ArgumentError


class VacancyHandlers(VacancyServicer):
    def __init__(self, vacancy_service: IVacancyService) -> None:
        self.vacancy_service: IVacancyService = vacancy_service
        super().__init__()

    async def CreateVacancy(self, request: CreateVacancyRequest, context: ServicerContext) -> CreateVacancyResponse:
        try:
            vacancy_id = await self.vacancy_service.create_vacancy(request)
            return CreateVacancyResponse(vacancy_id=vacancy_id)

        except ArgumentError as e:
            logging.info(f'CreateVacancy: {e}')

            context.abort(StatusCode.INVALID_ARGUMENT, e.__str__())

        except InternalError as e:
            logging.critical(f'CreateVacancy: {e}')

            context.abort(StatusCode.INTERNAL, e.__str__())
