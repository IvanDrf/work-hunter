from typing import Protocol

from src.domain.models.vacancy import VacancyORM


class IVacancyRepo(Protocol):
    async def create_vacancy(self, vacancy: VacancyORM) -> int:
        ...

    async def find_vacancy_by_id(self, vacancy_id: int) -> VacancyORM | None:
        ...
