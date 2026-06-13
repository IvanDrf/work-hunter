from asyncio import create_task
from datetime import timedelta

from src.core.exc import AccessError, ArgumentError
from src.domain.rules.user import is_user_admin
from src.domain.rules.vacancy import has_right_to_vacancy, is_vacancy_id_valid
from src.domain.schemas import UserInfo, VacancyResponseSchema
from src.infrastructure.service.base_vacancy import BaseVacancyService
from src.infrastructure.service.dependencies import ICache, IUnitOfWork, IVacancyRepo
from src.infrastructure.service.dto.vacancy_dto import vacancy_orm_to_response_dto


class VacancySearchService(BaseVacancyService):
    def __init__(
        self,
        vacancy_repo: IVacancyRepo,
        uof: IUnitOfWork,
        cache: ICache,
        vacancy_ttl: timedelta,
        cache_timeout: float,
    ) -> None:
        BaseVacancyService.__init__(self, cache, vacancy_ttl, cache_timeout)

        self.vacancy_repo: IVacancyRepo = vacancy_repo
        self.uof_factory: IUnitOfWork = uof

    async def find_vacancy_by_id(self, vacancy_id: int, user_info: UserInfo | None) -> VacancyResponseSchema | None:
        if not is_vacancy_id_valid(vacancy_id):
            raise ArgumentError(f"vacancy_id must be non negative number, {vacancy_id=}")

        vacancy = await self._get_vacancy_from_cache_by_id(vacancy_id)
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
        if not tags:
            raise ArgumentError("invalid vacancy tags, tags can't be empty")

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
        if not author:
            raise ArgumentError("invalid author name, author name is empty")

        vacancies = await self._get_vacancies_by_author(author, offset, limit)
        if vacancies is not None:
            return vacancies

        async with self.uof_factory as uof:
            if user_info is not None and is_user_admin(user_info):
                vacancies = await self.vacancy_repo.find_vacancies_for_admin_by_author(uof, author, offset, limit)
            else:
                vacancies = await self.vacancy_repo.find_only_published_vacancies_by_author(uof, author, offset, limit)

            if vacancies is None:
                return None

        vacancies = [vacancy_orm_to_response_dto(vacancy) for vacancy in vacancies]
        create_task(self._save_vacancies_by_author(vacancies, author, offset, limit))

        return vacancies

    async def find_vacancies_by_title(
        self,
        title: str,
        offset: int,
        limit: int,
        user_info: UserInfo | None,
    ) -> list[VacancyResponseSchema] | None:
        if not title:
            raise ArgumentError("invalid vacancy title, title can't be empty")

        async with self.uof_factory as uof:
            if user_info is not None and is_user_admin(user_info):
                vacancies = await self.vacancy_repo.find_vacancies_for_admin_by_title(uof, title, offset, limit)
            else:
                vacancies = await self.vacancy_repo.find_only_published_vacancies_by_title(uof, title, offset, limit)

        if vacancies is None:
            return None

        return [vacancy_orm_to_response_dto(vacancy) for vacancy in vacancies]
