from datetime import timedelta
from typing import Protocol, Self
from uuid import UUID

from src.domain.models import TagORM
from src.domain.models.vacancy import VacancyORM
from src.domain.types.enums import OrderBy


class ICache(Protocol):
    async def save(self, key: str, content: str, ttl: timedelta) -> None: ...
    async def get(self, key: str) -> str | None: ...


class IValidationServiceClient(Protocol):
    async def is_metro_valid(self, city: str, metro: str) -> bool: ...
    async def is_city_valid(self, city: str) -> bool: ...


class IUnitOfWork(Protocol):
    async def __aenter__(self) -> Self: ...
    async def __aexit__(self, exc_type, exc, tb): ...

    async def commit(self) -> None: ...
    async def rollback(self) -> None: ...
    async def flush(self) -> None: ...

    async def stop(self) -> None: ...


class IVacancySearchRepo(Protocol):
    async def find_vacancy_by_id(self, uof: IUnitOfWork, vacancy_id: int) -> VacancyORM | None: ...
    async def find_only_published_vacancies_with_tags(
        self,
        uof: IUnitOfWork,
        tags: list[str],
        order_by: OrderBy,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None: ...
    async def find_vacancies_for_admin_with_tags(
        self,
        uof: IUnitOfWork,
        tags: list[str],
        order_by: OrderBy,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None: ...

    async def find_only_published_vacancies_by_author(
        self,
        uof: IUnitOfWork,
        author: str,
        order_by: OrderBy,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None: ...
    async def find_vacancies_for_admin_by_author(
        self,
        uof: IUnitOfWork,
        author: str,
        order_by: OrderBy,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None: ...
    async def find_vacancies_by_author_id(
        self,
        uof: IUnitOfWork,
        author_id: UUID,
        order_by: OrderBy,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None: ...

    async def find_vacancies_for_admin_by_title(
        self,
        uof: IUnitOfWork,
        title: str,
        order_by: OrderBy,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None: ...
    async def find_only_published_vacancies_by_title(
        self,
        uof: IUnitOfWork,
        title: str,
        order_by: OrderBy,
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
