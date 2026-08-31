import logging
from asyncio import create_task, gather
from datetime import timedelta

from src.core.exc import AccessError, ArgumentError, InternalError
from src.domain.rules.user import is_user_admin, is_user_employee, is_user_employer
from src.domain.rules.vacancy import has_right_to_vacancy, is_vacancy_id_valid
from src.domain.schemas import UserInfo, VacancyResponseSchema
from src.domain.types import OrderBy, VacancyStatus
from src.infrastructure.service.dependencies import ICache, IUnitOfWork, IVacancyRepo
from src.infrastructure.service.vacancy.base_vacancy import BaseVacancyService
from src.infrastructure.service.vacancy.vacancy_dto import vacancy_orm_to_response_dto


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
        self.uof: IUnitOfWork = uof

        self.logger = logging.getLogger("VacancySearchService")

    async def find_vacancy_by_id(self, vacancy_id: int, user_info: UserInfo | None) -> VacancyResponseSchema | None:
        if not is_vacancy_id_valid(vacancy_id):
            raise ArgumentError(f"vacancy_id must be non negative number, {vacancy_id=}")

        vacancy = await self._get_vacancy_from_cache_by_id(vacancy_id)
        if vacancy is not None and has_right_to_vacancy(vacancy, user_info):
            create_task(self._update_vacancy_views(vacancy, user_info))
            return vacancy

        async with self.uof as uof:
            vacancy = await self.vacancy_repo.find_vacancy_by_id(uof, vacancy_id)
            if vacancy is None:
                return None

        if not has_right_to_vacancy(vacancy, user_info):
            raise AccessError("this vacancy is moderating now or deleted, you can't see it now")

        vacancy_schema = vacancy_orm_to_response_dto(vacancy)
        gather(*[self._save_vacancy_in_cache_by_id(vacancy_schema), self._update_vacancy_views(vacancy_schema, user_info)])

        return vacancy_schema

    async def find_vacancies_with_tags(
        self,
        tags: list[str],
        offset: int,
        limit: int,
        order_by: OrderBy,
        user_info: UserInfo | None,
    ) -> list[VacancyResponseSchema] | None:
        if not tags:
            raise ArgumentError("invalid vacancy tags, tags can't be empty")

        async with self.uof as uof:
            if user_info is not None and is_user_admin(user_info):
                vacancies = await self.vacancy_repo.find_vacancies_for_admin_with_tags(uof, tags, order_by, offset, limit)
            else:
                vacancies = await self.vacancy_repo.find_only_published_vacancies_with_tags(uof, tags, order_by, offset, limit)

            if vacancies is None:
                return None

        return [vacancy_orm_to_response_dto(vacancy) for vacancy in vacancies]

    async def find_vacancies_by_author(
        self,
        author: str,
        offset: int,
        limit: int,
        order_by: OrderBy,
        user_info: UserInfo | None,
    ) -> list[VacancyResponseSchema] | None:
        if not author:
            raise ArgumentError("invalid author name, author name is empty")

        vacancies = await self._get_vacancies_by_author(author, order_by, user_info, offset, limit)
        if vacancies is not None:
            return vacancies

        async with self.uof as uof:
            if user_info is not None and is_user_admin(user_info):
                vacancies = await self.vacancy_repo.find_vacancies_for_admin_by_author(uof, author, order_by, offset, limit)
            else:
                vacancies = await self.vacancy_repo.find_only_published_vacancies_by_author(uof, author, order_by, offset, limit)

            if vacancies is None:
                return None

        vacancies = [vacancy_orm_to_response_dto(vacancy) for vacancy in vacancies]
        create_task(self._save_vacancies_by_author(vacancies, author, order_by, user_info, offset, limit))

        return vacancies

    async def find_vacancies_by_author_id(
        self,
        offset: int,
        limit: int,
        order_by: OrderBy,
        user_info: UserInfo,
    ) -> list[VacancyResponseSchema] | None:
        if not is_user_employer(user_info) and not is_user_admin(user_info):
            raise AccessError(f"only employer or admin can see his vacancies, but user_role={user_info.user_role.name}")

        async with self.uof as uof:
            vacancies = await self.vacancy_repo.find_vacancies_by_author_id(uof, user_info.user_id, order_by, offset, limit)

        if vacancies is None:
            return None

        return [vacancy_orm_to_response_dto(vacancy) for vacancy in vacancies]

    async def find_vacancies_by_title(
        self,
        title: str,
        offset: int,
        limit: int,
        order_by: OrderBy,
        user_info: UserInfo | None,
    ) -> list[VacancyResponseSchema] | None:
        if not title:
            raise ArgumentError("invalid vacancy title, title can't be empty")

        async with self.uof as uof:
            if user_info is not None and is_user_admin(user_info):
                vacancies = await self.vacancy_repo.find_vacancies_for_admin_by_title(uof, title, order_by, offset, limit)
            else:
                vacancies = await self.vacancy_repo.find_only_published_vacancies_by_title(uof, title, order_by, offset, limit)

        if vacancies is None:
            return None

        return [vacancy_orm_to_response_dto(vacancy) for vacancy in vacancies]

    async def _update_vacancy_views(self, vacancy: VacancyResponseSchema, user_info: UserInfo | None) -> None:
        if vacancy.status != VacancyStatus.PUBLISHED:
            return

        if user_info is not None and (not is_user_employee(user_info) or user_info.user_id == vacancy.author_id):
            return

        async with self.uof as uof:
            try:
                await self.vacancy_repo.update_vacancy(uof, vacancy.vacancy_id, {"views": vacancy.views + 1})
            except InternalError as e:
                self.logger.critical(f"can't update views for {vacancy.vacancy_id=}, details={e}")
