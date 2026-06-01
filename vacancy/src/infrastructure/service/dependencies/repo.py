from typing import Protocol
from uuid import UUID

from src.domain.models import TagORM
from src.domain.models.vacancy import VacancyORM
from src.infrastructure.service.dependencies.unit_of_work import IUnitOfWork


class IVacancySearchRepo(Protocol):
    async def find_vacancy_by_id(self, uof: IUnitOfWork, vacancy_id: int) -> VacancyORM | None: ...
    async def find_only_published_vacancies_with_tags(
        self,
        uof: IUnitOfWork,
        tags: list[str],
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None: ...
    async def find_vacancies_for_admin_with_tags(
        self,
        uof: IUnitOfWork,
        tags: list[str],
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None: ...
    async def find_only_published_vacancies_by_author(
        self,
        uof: IUnitOfWork,
        author: str,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None: ...
    async def find_vacancies_for_admin_by_author(
        self,
        uof: IUnitOfWork,
        author: str,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None: ...
    async def find_vacancies_for_admin_by_title(
        self,
        uof: IUnitOfWork,
        title: str,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None: ...
    async def find_only_published_vacancies_by_title(
        self,
        uof: IUnitOfWork,
        title: str,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None: ...
    async def find_vacancy_author(self, uof: IUnitOfWork, vacancy_id: int) -> UUID | None: ...


class IVacancyDMLRepo(Protocol):
    async def create_vacancy(self, uof: IUnitOfWork, vacancy: VacancyORM, tags: list[TagORM]) -> None: ...
    async def update_vacancy(self, uof: IUnitOfWork, vacancy_id: int, fields: dict) -> VacancyORM: ...
    async def delete_vacancy(self, uof: IUnitOfWork, vacancy_id: int) -> None: ...


class IVacancyRepo(IVacancySearchRepo, IVacancyDMLRepo):
    pass


class ITagRepo(Protocol):
    async def add_tags(self, uof: IUnitOfWork, tags: list[str]) -> list[TagORM]: ...
