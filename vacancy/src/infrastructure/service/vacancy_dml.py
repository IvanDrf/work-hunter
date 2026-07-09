from asyncio import create_task, gather
from datetime import datetime, timedelta, timezone

from src.core.exc import AccessError, ArgumentError, NotFoundError
from src.domain.models.vacancy import VacancyORM, VacancyStatus
from src.domain.rules.user import is_user_admin, is_user_employer, is_user_vacancy_author
from src.domain.rules.vacancy import is_vacancy_id_valid
from src.domain.schemas import UserInfo, VacancyCreateSchema, VacancyResponseSchema, VacancyUpdateSchema
from src.domain.types import UNSET_VALUE, UnsetValue
from src.infrastructure.service.base_vacancy import BaseVacancyService
from src.infrastructure.service.dependencies import ICache, ITagRepo, IUnitOfWork, IVacancyRepo, IValidationServiceClient
from src.infrastructure.service.dto.vacancy_dto import (
    create_vacancy_dto,
    vacancy_orm_to_response_dto,
)


class VacancyDMLService(BaseVacancyService):
    def __init__(
        self,
        vacancy_repo: IVacancyRepo,
        tag_repo: ITagRepo,
        uof: IUnitOfWork,
        cache: ICache,
        validation_client: IValidationServiceClient,
        vacancy_ttl: timedelta,
        cache_timeout: float,
    ) -> None:
        BaseVacancyService.__init__(self, cache, vacancy_ttl, cache_timeout)

        self.vacancy_repo: IVacancyRepo = vacancy_repo
        self.tag_repo: ITagRepo = tag_repo
        self.uof_factory: IUnitOfWork = uof

        self.validation_client: IValidationServiceClient = validation_client

    async def create_vacancy(self, vacancy: VacancyCreateSchema, user_info: UserInfo) -> VacancyResponseSchema:
        if not vacancy.author_name:
            raise ArgumentError("author name can't be empty")

        if vacancy.city is None and vacancy.metro is not None:
            raise ArgumentError(f"city is empty, but metro doesn't, metro={vacancy.metro}")

        if not user_info.verificated:
            raise AccessError("user is not verificated, can't create vacancy")

        if not is_user_employer(user_info):
            raise AccessError("only employer can create vacancies")

        vacancy_create_date = datetime.now(timezone.utc)
        vacancy_status = VacancyStatus.MODERATING
        vacancyORM = create_vacancy_dto(vacancy, user_info, vacancy_create_date, vacancy_status)

        vacancyORM.is_city_valid, vacancyORM.is_metro_valid = await self.__validate_city_and_metro(vacancy)

        async with self.uof_factory as uof:
            tags = await self.tag_repo.add_tags(uof, vacancy.tags)
            await self.vacancy_repo.create_vacancy(uof, vacancyORM, tags)

        vacancy_schema = vacancy_orm_to_response_dto(vacancyORM)
        create_task(self._save_vacancy_in_cache_by_id(vacancy_schema))

        return vacancy_schema

    async def set_vacancy_status(
        self,
        vacancy_id: int,
        status: VacancyStatus,
        moderator_comments: str,
        user_info: UserInfo,
    ) -> None:
        if not is_vacancy_id_valid(vacancy_id):
            raise ArgumentError(f"vacancy_id must be non negative number, {vacancy_id=}")

        if not is_user_admin(user_info):
            raise AccessError("you can't change vacancy status, you are not admin")

        async with self.uof_factory as uof:
            fields = create_fields_for_status_update(status, moderator_comments)
            await self.vacancy_repo.update_vacancy(uof, vacancy_id, fields)

    async def delete_vacancy(self, vacancy_id: int, user_info: UserInfo) -> None:
        if not is_vacancy_id_valid(vacancy_id):
            raise ArgumentError(f"vacancy_id must be non negative number, {vacancy_id=}")

        if is_user_admin(user_info):
            await self.__delete_vacancy_by_admin(vacancy_id)
            return

        async with self.uof_factory as uof:
            author_id = await self.vacancy_repo.find_vacancy_author(uof, vacancy_id)

        if author_id is None:
            raise ArgumentError(f"can't find author for vacancy with {vacancy_id=}")

        if not is_user_vacancy_author(author_id, user_info.user_id):
            raise AccessError("you have no rights to delete vacancy, you didn't created")

        async with self.uof_factory as uof:
            await self.vacancy_repo.delete_vacancy(uof, vacancy_id)

    async def update_vacancy(self, vacancy_update_schema: VacancyUpdateSchema, user_info: UserInfo) -> VacancyResponseSchema:
        if not is_vacancy_id_valid(vacancy_update_schema.vacancy_id):
            raise ArgumentError(f"vacancy_id must be non negative number, {vacancy_update_schema.vacancy_id=}")

        fields = create_fields_for_update(vacancy_update_schema)

        async with self.uof_factory as uof:
            vacancy = await self.vacancy_repo.find_vacancy_by_id(uof, vacancy_update_schema.vacancy_id)
            if vacancy is None:
                raise NotFoundError(f"can't find vacancy with given vacancy_id={vacancy_update_schema.vacancy_id}")

            if not is_user_admin(user_info) and not is_user_vacancy_author(vacancy.author_id, user_info.user_id):
                raise AccessError("only author or admin can change vacancy")

            if not isinstance(vacancy_update_schema.tags, UnsetValue):
                tags = await self.tag_repo.add_tags(uof, vacancy_update_schema.tags)

                if not tags and not fields:
                    raise ArgumentError("no new fields to update vacancy")

                vacancy.tags = tags

                await uof.flush()

            if not fields:
                return vacancy_orm_to_response_dto(vacancy)

            await self.__update_vacancy_is_metro_valid(vacancy, fields)

            vacancy = await self.vacancy_repo.update_vacancy(uof, vacancy.vacancy_id, fields)
        return vacancy_orm_to_response_dto(vacancy)

    async def __delete_vacancy_by_admin(self, vacancy_id: int) -> None:
        async with self.uof_factory as uof:
            vacancy = await self.vacancy_repo.find_vacancy_by_id(uof, vacancy_id)
            if vacancy is None:
                raise NotFoundError(f"can't find vacancy with {vacancy_id=}")

            await self.vacancy_repo.delete_vacancy(uof, vacancy_id)

    async def __update_vacancy_is_metro_valid(self, vacancy: VacancyORM, fields: dict) -> None:
        city = fields.get("city", UNSET_VALUE)
        if city is UNSET_VALUE:
            city = vacancy.city

        metro = fields.get("metro", UNSET_VALUE)
        if metro is UNSET_VALUE:
            metro = vacancy.metro

        if city is None or metro is None:
            vacancy.is_metro_valid = False
            return

        if city is not UNSET_VALUE and metro is not UNSET_VALUE:
            vacancy.is_metro_valid = await self.validation_client.is_metro_valid(city, metro)

    async def __validate_city_and_metro(self, vacancy: VacancyCreateSchema) -> tuple[bool, bool]:
        is_city_valid = False
        is_metro_valid = False

        if vacancy.city is not None and vacancy.metro is not None:
            is_city_valid, is_metro_valid = await gather(
                *[
                    self.validation_client.is_city_valid(vacancy.city),
                    self.validation_client.is_metro_valid(vacancy.city, vacancy.metro),
                ]
            )
        elif vacancy.city is not None:
            is_city_valid = await self.validation_client.is_city_valid(vacancy.city)
            return is_city_valid, is_metro_valid

        return is_city_valid, is_metro_valid


def create_fields_for_update(vacancy_update_schema: VacancyUpdateSchema) -> dict:
    fields = {}
    for field_name, value in vacancy_update_schema.model_dump().items():
        if value is not UNSET_VALUE and field_name not in ("vacancy_id", "tags"):
            fields[field_name] = value

    if not fields and isinstance(vacancy_update_schema.tags, UnsetValue):
        raise ArgumentError("no fields were given to update")

    fields["updated_at"] = datetime.now(timezone.utc)
    return fields


def create_fields_for_status_update(status: VacancyStatus, moderator_comments: str) -> dict:
    fields: dict = {"status": status}
    updated_time = datetime.now(timezone.utc)

    match status:
        case VacancyStatus.PUBLISHED:
            fields["moderated_at"] = updated_time
            fields["moderator_comments"] = moderator_comments

        case VacancyStatus.CLOSED:
            fields["closed_at"] = updated_time
            fields["moderator_comments"] = moderator_comments if moderator_comments else None

        case VacancyStatus.DELETED:
            fields["deleted_at"] = updated_time
            fields["moderator_comments"] = moderator_comments

    return fields
