import logging
from asyncio import create_task, wait_for
from datetime import datetime, timedelta, timezone
from typing import Final

from pydantic import ValidationError

from src.core.exc import AccessError, ArgumentError, InternalError, NotFoundError
from src.domain.models.vacancy import VacancyStatus
from src.domain.rules.user import is_user_admin, is_user_employer, is_user_vacancy_author
from src.domain.rules.vacancy import has_right_to_vacancy, is_vacancy_id_valid
from src.domain.schemas import UserInfo, VacancyCreateSchema, VacancyResponseSchema, VacancyUpdateSchema
from src.domain.types.types import UNSET_VALUE, UnsetValue
from src.infrastructure.service.dependencies import ICache, ITagRepo, IUnitOfWork, IVacancyRepo
from src.infrastructure.service.dto.vacancy_dto import (
    create_vacancy_dto,
    vacancy_orm_to_response_dto,
)


class VacancyService:
    def __init__(
        self,
        vacancy_repo: IVacancyRepo,
        tag_repo: ITagRepo,
        cache: ICache,
        vacancy_ttl: float,
        unit_of_work: IUnitOfWork,
    ) -> None:
        self.vacancy_repo: IVacancyRepo = vacancy_repo
        self.tag_repo: ITagRepo = tag_repo

        self.cache: ICache = cache
        self.vacancy_ttl: timedelta = timedelta(minutes=vacancy_ttl)

        self.uof_factory: IUnitOfWork = unit_of_work
        self.CACHE_TIMEOUT: Final[float] = 0.3

    async def create_vacancy(self, vacancy: VacancyCreateSchema, user_info: UserInfo) -> VacancyResponseSchema:
        if not vacancy.author_name:
            raise ArgumentError("author name can't be empty")

        if not user_info.verificated:
            raise AccessError("user is not verificated, can't create vacancy")

        if not is_user_employer(user_info):
            raise AccessError("only employer can create vacancies")

        vacancy_create_date = datetime.now(timezone.utc)
        vacancy_status = VacancyStatus.MODERATING
        vacancyORM = create_vacancy_dto(vacancy, user_info, vacancy_create_date, vacancy_status)

        async with self.uof_factory as uof:
            tags = await self.tag_repo.add_tags(uof, vacancy.tags)
            await self.vacancy_repo.create_vacancy(uof, vacancyORM, tags)

        vacancySchema = vacancy_orm_to_response_dto(vacancyORM)
        create_task(self._save_vacancy_in_cache_by_id(vacancySchema))

        return vacancySchema

    async def find_vacancy_by_id(self, vacancy_id: int, user_info: UserInfo | None) -> VacancyResponseSchema | None:
        if not is_vacancy_id_valid(vacancy_id):
            raise ArgumentError(f"vacancy_id must be non negative number, {vacancy_id=}")

        vacancy = await self._get_vacancy_from_cache(vacancy_id)
        if vacancy is not None:
            if not has_right_to_vacancy(vacancy, user_info):
                raise AccessError("this vacancy is moderating now or deleted, you can't see it now")

            return vacancy

        async with self.uof_factory as uof:
            vacancy = await self.vacancy_repo.find_vacancy_by_id(uof, vacancy_id)
            if vacancy is None:
                return None

        if not has_right_to_vacancy(vacancy, user_info):
            raise AccessError("this vacancy is moderating now or deleted, you can't see it now")

        vacancySchema = vacancy_orm_to_response_dto(vacancy)
        create_task(self._save_vacancy_in_cache_by_id(vacancySchema))

        return vacancySchema

    async def find_vacancies_with_tags(
        self,
        tags: list[str],
        offset: int,
        limit: int,
        user_info: UserInfo | None,
    ) -> list[VacancyResponseSchema] | None:
        async with self.uof_factory as uof:
            if user_info is not None and is_user_admin(user_info):
                vacancies = await self.vacancy_repo.find_vacancies_for_admin_with_tags(uof, tags, offset, limit)
            else:
                vacancies = await self.vacancy_repo.find_only_published_vacancies_with_tags(uof, tags, offset, limit)

            if vacancies is None:
                return None

        return [vacancy_orm_to_response_dto(vacancy) for vacancy in vacancies]

    async def find_vacancies_by_author(
        self,
        author: str,
        offset: int,
        limit: int,
        user_info: UserInfo | None,
    ) -> list[VacancyResponseSchema] | None:
        if author == "":
            raise ArgumentError("invalid author name in request, author name is empty")

        async with self.uof_factory as uof:
            if user_info is not None and is_user_admin(user_info):
                vacancies = await self.vacancy_repo.find_vacancies_for_admin_by_author(uof, author, offset, limit)
            else:
                vacancies = await self.vacancy_repo.find_only_published_vacancies_by_author(uof, author, offset, limit)

            if vacancies is None:
                return None

        return [vacancy_orm_to_response_dto(vacancy) for vacancy in vacancies]

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

            vacancy = await self.vacancy_repo.update_vacancy(uof, vacancy.vacancy_id, fields)
        return vacancy_orm_to_response_dto(vacancy)

    async def __delete_vacancy_by_admin(self, vacancy_id: int) -> None:
        async with self.uof_factory as uof:
            vacancy = await self.vacancy_repo.find_vacancy_by_id(uof, vacancy_id)
            if vacancy is None:
                raise NotFoundError(f"can't find vacancy with {vacancy_id=}")

            await self.vacancy_repo.delete_vacancy(uof, vacancy_id)

    async def _save_vacancy_in_cache_by_id(self, vacancy: VacancyResponseSchema) -> None:
        vacancy_json = vacancy.model_dump_json()

        try:
            await self.cache.save(str(vacancy.vacancy_id), vacancy_json, self.vacancy_ttl)
        except InternalError as e:
            logging.warning(f"Can't save vacancy with vacancy_id={vacancy.vacancy_id} in cache, error={e}")

    async def _get_vacancy_from_cache(self, vacancy_id: int) -> VacancyResponseSchema | None:
        try:
            vacancy_json = await wait_for(self.cache.get(str(vacancy_id)), timeout=self.CACHE_TIMEOUT)
            return VacancyResponseSchema.model_validate_json(vacancy_json) if vacancy_json is not None else None

        except InternalError as e:
            logging.warning(f"Can't get vacancy with {vacancy_id=} from cache, error={e}")

        except ValidationError as e:
            logging.error(f"Invalid vacancy json in cache, can't parse it, error={e}")

        except TimeoutError:
            logging.error(f"Can't get vacancy with {vacancy_id=} form cache, timeout error")


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
