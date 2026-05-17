from grpc import ServicerContext
from pkg.vacancy_api.vacancy_pb2 import (
    CreateVacancyRequest,
    DeleteVacancyRequest,
    FindVacancyByIDRequest,
    FindVacancyByTagsRequest,
    Response,
    ResponseStatus,
    SetVacancyStatusRequest,
    UpdateVacancyRequest,
    Vacancies,
    VacancyInfo,
)
from pkg.vacancy_api.vacancy_pb2_grpc import VacancyServicer

from src.api.dependencies import IVacancyService
from src.api.dto.status import vacancy_status_dto
from src.api.dto.user_info import user_info_dto, user_info_none_dto
from src.api.dto.vacancy import vacancy_create_dto, vacancy_response_dto, vacancy_update_dto
from src.api.rules.params import (
    MAX_LIMIT,
    MAX_TAGS_AMOUNT,
    MIN_LIMIT,
    MIN_OFFSET,
    MIN_TAGS_AMOUNT,
    is_limit_valid,
    is_offset_valid,
    is_tags_amount_valid,
)
from src.api.rules.user_info import get_user_info, is_user_id_valid
from src.core.exc import ArgumentError, NotFoundError
from src.utils.handle_errors import handle_errors


class VacancyHandlers(VacancyServicer):
    def __init__(self, vacancy_service: IVacancyService) -> None:
        self.vacancy_service: IVacancyService = vacancy_service
        super().__init__()

    @handle_errors
    async def CreateVacancy(self, request: CreateVacancyRequest, context: ServicerContext) -> VacancyInfo:
        if not is_user_id_valid(request.user_info):
            raise ArgumentError(f"invalid user_id in user_info: {request.user_info.user_id=}")

        vacancy_schema = vacancy_create_dto(request.vacancy, request.user_info)
        user_info_schema = user_info_dto(request.user_info)

        vacancy = await self.vacancy_service.create_vacancy(vacancy_schema, user_info_schema)

        return vacancy_response_dto(vacancy)

    @handle_errors
    async def FindVacancyByID(self, request: FindVacancyByIDRequest, context: ServicerContext) -> VacancyInfo:
        user_info = get_user_info(request)

        vacancy = await self.vacancy_service.find_vacancy_by_id(request.vacancy_id, user_info_none_dto(user_info))
        if vacancy is None:
            raise NotFoundError(f"can't find vacancy with given {request.vacancy_id=}")

        return vacancy_response_dto(vacancy)

    @handle_errors
    async def FindVacanciesByTags(self, request: FindVacancyByTagsRequest, context: ServicerContext) -> Vacancies:
        if not is_offset_valid(request.offset):
            raise ArgumentError(f"offset must be greater than {MIN_OFFSET}, but {request.offset=}")

        if not is_limit_valid(request.limit):
            raise ArgumentError(f"limit must be in range ({MIN_LIMIT}, {MAX_LIMIT}), but {request.limit=}")

        if not is_tags_amount_valid(request.tags):
            raise ArgumentError(f"tags amount must be in range ({MIN_TAGS_AMOUNT}, {MAX_TAGS_AMOUNT}), but {len(request.tags)=}")

        user_info = get_user_info(request)

        vacancies = await self.vacancy_service.find_vacancies_with_tags(
            list(request.tags), request.offset, request.limit, user_info_none_dto(user_info)
        )

        if vacancies is None:
            raise NotFoundError(
                f"can't find vacancies with given params {list(request.tags)}, {request.offset=}, {request.limit=}"
            )

        return Vacancies(
            vacancies=[vacancy_response_dto(vacancy) for vacancy in vacancies],
            limit=request.limit,
            offset=request.offset,
        )

    @handle_errors
    async def SetVacancyStatus(self, request: SetVacancyStatusRequest, context: ServicerContext) -> Response:
        await self.vacancy_service.set_vacancy_status(
            request.vacancy_id,
            vacancy_status_dto(request.status),
            request.moderator_comments,
            user_info_dto(request.user_info),
        )

        return Response(message="successfully updated vacancy status", status=ResponseStatus.SUCCESS)

    @handle_errors
    async def DeleteVacancy(self, request: DeleteVacancyRequest, context: ServicerContext) -> Response:
        await self.vacancy_service.delete_vacancy(request.vacancy_id, user_info_dto(request.user_info))

        return Response(message="successfully deleted vacancy", status=ResponseStatus.SUCCESS)

    @handle_errors
    async def UpdateVacancy(self, request: UpdateVacancyRequest, context: ServicerContext) -> VacancyInfo:
        user_info = user_info_dto(request.user_info)
        vacancy_schema = vacancy_update_dto(request)

        vacancy = await self.vacancy_service.update_vacancy(vacancy_schema, user_info)
        return vacancy_response_dto(vacancy)
