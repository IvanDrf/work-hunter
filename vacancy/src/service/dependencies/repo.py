from typing import Protocol

from src.domain.models.vacancy import VacancyORM, VacancyStatus


class IVacancyRepo(Protocol):
    async def create_vacancy(self, vacancy: VacancyORM) -> int:
        ...

    async def find_vacancy_by_id(self, vacancy_id: int) -> VacancyORM | None:
        ...

    async def find_vacancies_with_tags(self, tags: list[str], offset: int, limit: int) -> list[VacancyORM] | None:
        ...

    async def set_vacancy_status(self, vacancy_id: int, status: VacancyStatus) -> None:
        ...
