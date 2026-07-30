from uuid import UUID

from sqlalchemy import and_, select
from sqlalchemy.exc import DBAPIError, SQLAlchemyError

from src.core.exc import InternalError
from src.domain.models import ApplicationORM
from src.infrastructure.persistence.postgresql_repo.uof import UnitOfWork
from src.infrastructure.persistence.postgresql_repo.utils import catch_and_raise

DB_ERRORS = (SQLAlchemyError, DBAPIError)


class ApplicationPostgreSQLRepo:
    @catch_and_raise(DB_ERRORS, InternalError, "can't add new application in database")
    async def add_application(self, uof: UnitOfWork, application: ApplicationORM) -> None:
        uof.session.add(application)

    @catch_and_raise(DB_ERRORS, InternalError, "can't find application with given params")
    async def find_application(self, uof: UnitOfWork, vacancy_id: int, user_id: UUID) -> ApplicationORM | None:
        query = select(ApplicationORM).where(and_(ApplicationORM.vacancy_id == vacancy_id, ApplicationORM.user_id == user_id))

        res = await uof.session.execute(query)
        return res.scalar_one_or_none()

    @catch_and_raise(DB_ERRORS, InternalError, "can't find vacancies ids with given params")
    async def find_vacancies_ids_by_user_id(self, uof: UnitOfWork, user_id: UUID, *, limit: int, offset: int) -> list[int] | None:
        query = select(ApplicationORM.vacancy_id).where(ApplicationORM.user_id == user_id).limit(limit).offset(offset)

        res = await uof.session.execute(query)
        rows = res.scalars().all()
        return list(rows) if rows else None
