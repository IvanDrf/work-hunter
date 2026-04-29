from grpc import ServicerContext
from pkg.vacancy_api.vacancy_pb2 import (CreateVacancyRequest, CreateVacancyResponse, DeleteVacancyRequest,
                                         FindVacancyByIDRequest, FindVacancyByTagsRequest, Response, ResponseStatus,
                                         SetVacancyStatusRequest, Vacancies, VacancyInfo,)
from pkg.vacancy_api.vacancy_pb2_grpc import VacancyServicer

from src.api.dependencies.service import IVacancyService
from src.api.rules.params import (MAX_LIMIT, MAX_TAGS_AMOUNT, MIN_LIMIT, MIN_OFFSET, MIN_TAGS_AMOUNT, is_limit_valid,
                                  is_offset_valid, is_tags_amount_valid,)
from src.core.exc import ArgumentError, NotFoundError
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

    @handle_errors
    async def FindVacanciesByTags(self, request: FindVacancyByTagsRequest, context: ServicerContext) -> Vacancies:
        if not is_offset_valid(request.offset):
            raise ArgumentError(
                f'offset must be greater than {MIN_OFFSET}, but {request.offset=}')

        if not is_limit_valid(request.limit):
            raise ArgumentError(
                f'limit must be in range ({MIN_LIMIT}, {MAX_LIMIT}), but {request.limit=}'
            )

        if not is_tags_amount_valid(request.tags):
            raise ArgumentError(
                f'tags amount must be in range ({MIN_TAGS_AMOUNT}, {MAX_TAGS_AMOUNT}), but {len(request.tags)=}'
            )

        vacancies = await self.vacancy_service.find_vacancies_with_tags(list(request.tags), request.offset, request.limit)
        if vacancies is None:
            raise NotFoundError(
                f'''can't find vacancies with given params {list(request.tags)}, {request.offset=}, {request.limit=}'''
            )

        return vacancies

    @handle_errors
    async def SetVacancyStatus(self, request: SetVacancyStatusRequest, context: ServicerContext) -> Response:
        await self.vacancy_service.set_vacancy_status(request.vacancy_id, request.status, request.moderator_comments, request.user_info)

        return Response(message='successfully updated vacancy status', status=ResponseStatus.SUCCESS)

    @handle_errors
    async def DeleteVacancy(self, request: DeleteVacancyRequest, context: ServicerContext) -> Response:
        await self.vacancy_service.delete_vacancy(request.vacancy_id, request.user_info)

        return Response(message='successfully deleted vacancy', status=ResponseStatus.SUCCESS)
