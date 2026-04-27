from typing import Protocol

from pkg.common.common_pb2 import UserInfo
from pkg.vacancy_api.vacancy_pb2 import UpdateVacancyRequest, Vacancies, VacancyInfo
from pkg.vacancy_api.vacancy_pb2 import VacancyStatus as PBVacancyStatus


class IVacancyService(Protocol):
    async def create_vacancy(self, vacancy: VacancyInfo, user_info: UserInfo) -> int:
        ...

    async def find_vacancy_by_id(self, vacancy_id: int, user_info: UserInfo) -> VacancyInfo | None:
        ...

    async def find_vacancies_by_author(self, author: str, user_info: UserInfo) -> Vacancies:
        ...

    async def update_vacancy(self, request: UpdateVacancyRequest) -> None:
        ...

    async def delete_vacancy(self, vacancy_id: int, user_info: UserInfo) -> None:
        ...

    async def set_vacancy_status(self, vacancy_id: int, status: PBVacancyStatus, user_info: UserInfo) -> None:
        ...
