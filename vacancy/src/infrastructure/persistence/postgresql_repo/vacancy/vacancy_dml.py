from sqlalchemy import case, delete, update

from src.core.exc import InternalError
from src.domain.models import TagORM, VacancyORM
from src.domain.types import DBErrors
from src.infrastructure.persistence.postgresql_repo.unit_of_work import UnitOfWork
from src.utils.catch_error import catch_raise_error


class VacancyDMLRepo:
    @catch_raise_error(DBErrors, raise_error=InternalError, message="can't create new vacancy in database")
    async def create_vacancy(self, uof: UnitOfWork, vacancy: VacancyORM, tags: list[TagORM]) -> None:
        vacancy.tags = tags
        uof.session.add(vacancy)

    @catch_raise_error(DBErrors, raise_error=InternalError, message="can't delete vacancy with given vacancy_id")
    async def delete_vacancy(self, uof: UnitOfWork, vacancy_id: int) -> None:
        query = delete(VacancyORM).where(VacancyORM.vacancy_id == vacancy_id).returning(VacancyORM.vacancy_id)

        res = await uof.session.execute(query)
        if res.one_or_none() is None:
            raise InternalError(f"can't delete vacancy with {vacancy_id=}")

    @catch_raise_error(DBErrors, raise_error=InternalError, message="can't update vacancy with given vacancy_id")
    async def update_vacancy(self, uof: UnitOfWork, vacancy_id: int, fields: dict) -> VacancyORM:
        query = update(VacancyORM).where(VacancyORM.vacancy_id == vacancy_id).values(fields).returning(VacancyORM)

        res = await uof.session.execute(query)
        vacancy = res.scalar_one_or_none()

        if vacancy is None:
            raise InternalError(f"can't update vacancy with {vacancy_id=}")

        return vacancy

    @catch_raise_error(DBErrors, raise_error=InternalError, message="can't update vacancies")
    async def update_vacancies(self, uof: UnitOfWork, fields: list[dict]) -> None:
        query = (
            update(VacancyORM)
            .where(VacancyORM.vacancy_id.in_(field["vacancy_id"] for field in fields))
            .values(
                applications_count=VacancyORM.applications_count
                + case(*[VacancyORM.vacancy_id == field["vacancy_id"] for field in fields], else_=0)
            )
        )

        await uof.session.execute(query)
