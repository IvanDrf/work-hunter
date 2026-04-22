from typing import Protocol

from src.domain.models.vacancy import VacancyORM


class IVacancyRepo(Protocol):
    async def create_vacancy(self, vacancy: VacancyORM) -> int:
        ...
