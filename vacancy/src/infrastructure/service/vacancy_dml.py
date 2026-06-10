from asyncio import create_task
from datetime import datetime, timedelta, timezone

from src.core.exc import AccessError, ArgumentError, NotFoundError
from src.domain.models.vacancy import VacancyStatus
from src.domain.rules.user import is_user_admin, is_user_employer, is_user_vacancy_author
from src.domain.rules.vacancy import is_vacancy_id_valid
from src.domain.schemas import UserInfo, VacancyCreateSchema, VacancyResponseSchema, VacancyUpdateSchema
from src.domain.types.types import UNSET_VALUE, UnsetValue
from src.infrastructure.service.base_vacancy import BaseVacancyService
from src.infrastructure.service.dependencies import ICache, ITagRepo, IUnitOfWork, IVacancyRepo
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
        vacancy_ttl: timedelta,
        cache_timeout: float,
    ) -> None:
        BaseVacancyService.__init__(self, cache, vacancy_ttl, cache_timeout)

        self.vacancy_repo: IVacancyRepo = vacancy_repo
        self.tag_repo: ITagRepo = tag_repo

        self.uof_factory: IUnitOfWork = uof

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
