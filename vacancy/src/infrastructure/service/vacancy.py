from datetime import datetime, timezone

from src.core.exc import AccessError, ArgumentError, NotFoundError
from src.domain.models.vacancy import VacancyStatus
from src.domain.rules.user import is_user_admin, is_user_employer, is_user_vacancy_author
from src.domain.rules.vacancy import has_right_to_vacancy, is_vacancy_id_valid
from src.domain.schemas import UserInfo, VacancyCreateSchema, VacancyResponseSchema, VacancyUpdateSchema
from src.domain.types.types import UNSET_VALUE, UnsetValue
from src.infrastructure.service.dependencies import ITagRepo, IUnitOfWork, IVacancyRepo
from src.infrastructure.service.dto.vacancy import (
    create_vacancy_dto,
    vacancy_orm_to_response_dto,
)


class VacancyService:
    def __init__(self, vacancy_repo: IVacancyRepo, tag_repo: ITagRepo, unit_of_work: IUnitOfWork) -> None:
        self.vacancy_repo: IVacancyRepo = vacancy_repo
        self.tag_repo: ITagRepo = tag_repo
        self.uof_factory: IUnitOfWork = unit_of_work

    async def create_vacancy(self, vacancy: VacancyCreateSchema, user_info: UserInfo) -> VacancyResponseSchema:
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

        return vacancy_orm_to_response_dto(vacancyORM)

    async def find_vacancy_by_id(self, vacancy_id: int, user_info: UserInfo | None) -> VacancyResponseSchema | None:
        if not is_vacancy_id_valid(vacancy_id):
            raise ArgumentError(f"vacancy_id must be non negative number, {vacancy_id=}")

        async with self.uof_factory as uof:
            vacancy = await self.vacancy_repo.find_vacancy_by_id(uof, vacancy_id)
            if vacancy is None:
                return None

        if not has_right_to_vacancy(vacancy, user_info):
            raise AccessError("this vacancy is moderating now, you can't see it now")

        return vacancy_orm_to_response_dto(vacancy)

    async def find_vacancies_with_tags(
        self, tags: list[str], offset: int, limit: int, user_info: UserInfo | None
    ) -> list[VacancyResponseSchema] | None:
        async with self.uof_factory as uof:
            if user_info is not None and is_user_admin(user_info):
                vacancies = await self.vacancy_repo.find_vacancies_for_admin_with_tags(uof, tags, offset, limit)
            else:
                vacancies = await self.vacancy_repo.find_only_published_vacancies_with_tags(uof, tags, offset, limit)

            if vacancies is None:
                return None

        return [vacancy_orm_to_response_dto(vacancy) for vacancy in vacancies]

    async def set_vacancy_status(
        self, vacancy_id: int, status: VacancyStatus, moderator_comments: str, user_info: UserInfo
    ) -> None:
        if not is_vacancy_id_valid(vacancy_id):
            raise ArgumentError(f"vacancy_id must be non negative number, {vacancy_id=}")

        if not is_user_admin(user_info):
            raise AccessError("you can't change vacancy status, you are not admin")

        async with self.uof_factory as uof:
            moderated_at = datetime.now(timezone.utc)
            await self.vacancy_repo.set_vacancy_status(uof, vacancy_id, VacancyStatus(status), moderator_comments, moderated_at)

    async def delete_vacancy(self, vacancy_id: int, user_info: UserInfo) -> None:
        if not is_vacancy_id_valid(vacancy_id):
            raise ArgumentError(f"vacancy_id must be non negative number, {vacancy_id=}")

        if is_user_admin(user_info):
            async with self.uof_factory as uof:
                await self.vacancy_repo.delete_vacancy(uof, vacancy_id)
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

        fields = parse_updated_fields(vacancy_update_schema)

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


def parse_updated_fields(vacancy_update_schema: VacancyUpdateSchema) -> dict:
    fields = {}
    for field_name, value in vacancy_update_schema.model_dump().items():
        if value is not UNSET_VALUE and field_name not in ("vacancy_id", "tags"):
            fields[field_name] = value

    if not fields and isinstance(vacancy_update_schema.tags, UnsetValue):
        raise ArgumentError("no fields were given to update")

    fields["updated_at"] = datetime.now(timezone.utc)
    return fields
