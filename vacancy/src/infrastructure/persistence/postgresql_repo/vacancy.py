from uuid import UUID

from sqlalchemy import and_, delete, select, update
from sqlalchemy.exc import DBAPIError, SQLAlchemyError
from sqlalchemy.orm import selectinload

from src.core.exc import InternalError
from src.domain.models import TagORM, VacanciesTagsORM, VacancyORM
from src.domain.models.vacancy import VacancyStatus
from src.infrastructure.persistence.postgresql_repo.unit_of_work import UnitOfWork
from src.utils.catch_error import catch_raise_error


class VacancyRepo:
    @catch_raise_error((SQLAlchemyError, DBAPIError), InternalError, "critical", "can't create new vacancy in database")
    async def create_vacancy(self, uof: UnitOfWork, vacancy: VacancyORM, tags: list[TagORM]) -> None:
        vacancy.tags = tags
        uof.session.add(vacancy)

    @catch_raise_error((SQLAlchemyError, DBAPIError), InternalError, "critical", "can't find vacancy with given vacancy_id")
    async def find_vacancy_by_id(self, uof: UnitOfWork, vacancy_id: int) -> VacancyORM | None:
        query = select(VacancyORM).where(VacancyORM.vacancy_id == vacancy_id).options(selectinload(VacancyORM.tags))

        res = await uof.session.execute(query)
        return res.scalar_one_or_none()

    @catch_raise_error((SQLAlchemyError, DBAPIError), InternalError, "critical", "can't find vacancies with given tags")
    async def find_only_published_vacancies_with_tags(
        self,
        uof: UnitOfWork,
        tags: list[str],
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None:
        query = (
            select(VacancyORM)
            .join(VacanciesTagsORM, VacanciesTagsORM.vacancy_id == VacancyORM.vacancy_id)
            .join(TagORM, TagORM.tag_id == VacanciesTagsORM.tag_id)
            .where(and_(TagORM.tag.in_(tags), VacancyORM.status == VacancyStatus.PUBLISHED))
            .offset(offset)
            .limit(limit)
            .options(selectinload(VacancyORM.tags))
        )

        res = await uof.session.execute(query)
        vacancies = list(res.scalars())

        return vacancies if len(vacancies) > 0 else None

    @catch_raise_error((SQLAlchemyError, DBAPIError), InternalError, "critical", "can't find vacancies with given tags")
    async def find_vacancies_for_admin_with_tags(
        self,
        uof: UnitOfWork,
        tags: list[str],
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None:
        query = (
            select(VacancyORM)
            .join(VacanciesTagsORM, VacanciesTagsORM.vacancy_id == VacancyORM.vacancy_id)
            .join(TagORM, TagORM.tag_id == VacanciesTagsORM.tag_id)
            .where(TagORM.tag.in_(tags))
            .offset(offset)
            .limit(limit)
            .options(selectinload(VacancyORM.tags))
        )

        res = await uof.session.execute(query)
        vacancies = list(res.scalars())

        return vacancies if len(vacancies) > 0 else None

    @catch_raise_error((SQLAlchemyError, DBAPIError), InternalError, "critical", "can't find vacancies by author")
    async def find_only_published_vacancies_by_author(
        self,
        uof: UnitOfWork,
        author: str,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None:
        query = (
            select(VacancyORM)
            .where(and_(VacancyORM.status == VacancyStatus.PUBLISHED, VacancyORM.author_name.like(author)))
            .offset(offset)
            .limit(limit)
            .options(selectinload(VacancyORM.tags))
        )

        res = await uof.session.execute(query)
        vacancies = list(res.scalars().all())

        return vacancies if len(vacancies) > 0 else None

    @catch_raise_error((SQLAlchemyError, DBAPIError), InternalError, "critical", "can't find vacancies by author for admin")
    async def find_vacancies_for_admin_by_author(
        self,
        uof: UnitOfWork,
        author: str,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None:
        query = (
            select(VacancyORM)
            .where(VacancyORM.author_name.like(author))
            .offset(offset)
            .limit(limit)
            .options(selectinload(VacancyORM.tags))
        )

        res = await uof.session.execute(query)
        vacancies = list(res.scalars().all())

        return vacancies if len(vacancies) > 0 else None

    @catch_raise_error((SQLAlchemyError, DBAPIError), InternalError, "critical", "can't find author_id by vacancy_id")
    async def find_vacancy_author(self, uof: UnitOfWork, vacancy_id: int) -> UUID | None:
        query = select(VacancyORM.author_id).where(VacancyORM.vacancy_id == vacancy_id)

        author_id = await uof.session.execute(query)
        return author_id.scalar_one_or_none()

    @catch_raise_error((SQLAlchemyError, DBAPIError), InternalError, "critical", "can't delete vacancy with given vacancy_id")
    async def delete_vacancy(self, uof: UnitOfWork, vacancy_id: int) -> None:
        query = delete(VacancyORM).where(VacancyORM.vacancy_id == vacancy_id).returning(VacancyORM.vacancy_id)

        res = await uof.session.execute(query)
        if res.one_or_none() is None:
            raise InternalError(f"can't delete vacancy with {vacancy_id=}")

    @catch_raise_error((SQLAlchemyError, DBAPIError), InternalError, "critical", "can't update vacancy with given vacancy_id")
    async def update_vacancy(self, uof: UnitOfWork, vacancy_id: int, fields: dict) -> VacancyORM:
        query = update(VacancyORM).where(VacancyORM.vacancy_id == vacancy_id).values(fields).returning(VacancyORM)

        res = await uof.session.execute(query)
        vacancy = res.scalar_one_or_none()

        if vacancy is None:
            raise InternalError(f"can't update vacancy with {vacancy_id=}")

        return vacancy
