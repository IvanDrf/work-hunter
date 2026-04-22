from typing import Protocol

from pkg.vacancy_api.vacancy_pb2 import (CreateVacancyRequest, DeleteVacancyRequest, FindVacanciesByAuthorRequest,
                                         FindVacancyByIDRequest, SetVacancyStatusRequest, UpdateVacancyRequest,
                                         Vacancies, VacancyInfo,)


class IVacancyService(Protocol):
    async def create_vacancy(self, request: CreateVacancyRequest) -> int:
        ...

    async def find_vacancy_by_id(self, request: FindVacancyByIDRequest) -> VacancyInfo:
        ...

    async def find_vacancies_by_author(self, request: FindVacanciesByAuthorRequest) -> Vacancies:
        ...

    async def update_vacancy(self, request: UpdateVacancyRequest) -> None:
        ...

    async def delete_vacancy(self, request: DeleteVacancyRequest) -> None:
        ...

    async def set_vacancy_status(self, request: SetVacancyStatusRequest) -> None:
        ...
