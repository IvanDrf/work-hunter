from grpc import ServicerContext
from pkg.common.common_pb2 import Empty, ServiceStatus, Status
from pkg.vacancy_api.vacancy_pb2 import (
    CreateVacancyRequest,
    DeleteVacancyRequest,
    FindVacanciesByAuthorIDRequest,
    FindVacanciesByAuthorRequest,
    FindVacanciesByTitleRequest,
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

from src.api.grpc.dependencies import IVacancyService
from src.api.grpc.dto.order_by import order_by_dto
from src.api.grpc.dto.status import vacancy_status_dto
from src.api.grpc.dto.user_info import user_info_dto, user_info_none_dto
from src.api.grpc.dto.vacancy import vacancy_create_dto, vacancy_response_dto, vacancy_update_dto
from src.api.grpc.rules.params import (
    MAX_TAGS_AMOUNT,
    MIN_TAGS_AMOUNT,
    is_tags_amount_valid,
    is_tags_values_valid,
    validate_limit_offset,
)
from src.api.grpc.rules.user_info import get_user_info
from src.core.exc import AccessError, ArgumentError, NotFoundError
from src.utils.handle_errors import handle_errors


class VacancyHandlers(VacancyServicer):
    def __init__(self, vacancy_service: IVacancyService) -> None:
        self.vacancy_service: IVacancyService = vacancy_service
        super().__init__()

    async def stop(self) -> None:
        await self.vacancy_service.stop()

    async def Health(self, request: Empty, context: ServicerContext) -> ServiceStatus:
        return ServiceStatus(status=Status.AVAILABLE)

    @handle_errors
    async def CreateVacancy(self, request: CreateVacancyRequest, context: ServicerContext) -> VacancyInfo:
        if not is_tags_amount_valid(request.tags):
            raise ArgumentError(f"tags amount must be in range ({MIN_TAGS_AMOUNT}, {MAX_TAGS_AMOUNT}), but {len(request.tags)=}")

        if not is_tags_values_valid(request.tags):
            raise ArgumentError("empty tag is not allowed in tags")

        vacancy_schema = vacancy_create_dto(request)
        user_info_schema = user_info_dto(request.user_info)

        vacancy = await self.vacancy_service.create_vacancy(vacancy_schema, user_info_schema)

        return vacancy_response_dto(vacancy)

    @handle_errors
    async def FindVacancyByID(self, request: FindVacancyByIDRequest, context: ServicerContext) -> VacancyInfo:
        user_info = get_user_info(request)

        vacancy = await self.vacancy_service.find_vacancy_by_id(request.vacancy_id, user_info_none_dto(user_info))
        if vacancy is None:
            raise NotFoundError(f"can't find vacancy with vacancy_id={request.vacancy_id}")

        return vacancy_response_dto(vacancy)

    @handle_errors
    async def FindVacanciesByTags(self, request: FindVacancyByTagsRequest, context: ServicerContext) -> Vacancies:
        validate_limit_offset(request.limit, request.offset)

        if not is_tags_amount_valid(request.tags):
            raise ArgumentError(f"tags amount must be in range ({MIN_TAGS_AMOUNT}, {MAX_TAGS_AMOUNT}), but {len(request.tags)=}")

        if not is_tags_values_valid(request.tags):
            raise ArgumentError("empty tag is not allowed in tags")

        user_info = get_user_info(request)

        vacancies = await self.vacancy_service.find_vacancies_with_tags(
            tags=list(request.tags),
            offset=request.offset,
            limit=request.limit,
            order_by=order_by_dto(request),
            user_info=user_info_none_dto(user_info),
        )

        if vacancies is None:
            raise NotFoundError(
                f"can't find any vacancies with tags={list(request.tags)}, limit={request.limit}, offset={request.offset}"
            )

        return Vacancies(
            vacancies=[vacancy_response_dto(vacancy) for vacancy in vacancies],
            limit=request.limit,
            offset=request.offset,
        )

    @handle_errors
    async def SetVacancyStatus(self, request: SetVacancyStatusRequest, context: ServicerContext) -> Response:
        await self.vacancy_service.set_vacancy_status(
            vacancy_id=request.vacancy_id,
            status=vacancy_status_dto(request.status),
            moderator_comments=request.moderator_comments,
            user_info=user_info_dto(request.user_info),
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

        if not is_tags_amount_valid(vacancy_schema.tags):
            raise ArgumentError(f"tags amount must be in range ({MIN_TAGS_AMOUNT}, {MAX_TAGS_AMOUNT}), but {len(request.tags)=}")

        if not is_tags_values_valid(vacancy_schema.tags):
            raise ArgumentError("empty tag is not allowed in tags")

        vacancy = await self.vacancy_service.update_vacancy(vacancy_schema, user_info)
        return vacancy_response_dto(vacancy)

    @handle_errors
    async def FindVacanciesByAuthor(self, request: FindVacanciesByAuthorRequest, context: ServicerContext) -> Vacancies:
        validate_limit_offset(request.limit, request.offset)

        user_info = get_user_info(request)
        if user_info is not None:
            user_info = user_info_dto(user_info)

        vacancies = await self.vacancy_service.find_vacancies_by_author(
            author=request.author_name,
            offset=request.offset,
            limit=request.limit,
            order_by=order_by_dto(request),
            user_info=user_info,
        )
        if vacancies is None:
            raise NotFoundError(
                f"can't find any vacancies with author_name={request.author_name}, limit={request.limit}, offset={request.offset}"
            )

        return Vacancies(
            vacancies=[vacancy_response_dto(vacancy) for vacancy in vacancies],
            limit=request.limit,
            offset=request.offset,
        )

    @handle_errors
    async def FindVacanciesByTitle(self, request: FindVacanciesByTitleRequest, context: ServicerContext) -> Vacancies:
        validate_limit_offset(request.limit, request.offset)

        user_info = get_user_info(request)

        vacancies = await self.vacancy_service.find_vacancies_by_title(
            title=request.title,
            offset=request.offset,
            limit=request.limit,
            order_by=order_by_dto(request),
            user_info=user_info_none_dto(user_info),
        )

        if vacancies is None:
            raise NotFoundError(
                f"can't find any vacancies with title={request.title}, limit={request.limit}, offset={request.offset}"
            )

        return Vacancies(
            vacancies=[vacancy_response_dto(vacancy) for vacancy in vacancies],
            limit=request.limit,
            offset=request.offset,
        )

    @handle_errors
    async def FindVacanciesByAuthorID(self, request: FindVacanciesByAuthorIDRequest, context: ServicerContext) -> Vacancies:
        validate_limit_offset(request.limit, request.offset)

        user_info = get_user_info(request)
        if user_info is None:
            raise AccessError("invalid user_info in request, user_info is empty")

        vacancies = await self.vacancy_service.find_vacancies_by_author_id(
            offset=request.offset,
            limit=request.limit,
            order_by=order_by_dto(request),
            user_info=user_info_dto(user_info),
        )

        if vacancies is None:
            raise NotFoundError(
                f"can't find any vacancies with user_id={user_info.user_id} with limit={request.limit}, offset={request.offset}"
            )

        return Vacancies(
            vacancies=[vacancy_response_dto(vacancy) for vacancy in vacancies],
            limit=request.limit,
            offset=request.offset,
        )
