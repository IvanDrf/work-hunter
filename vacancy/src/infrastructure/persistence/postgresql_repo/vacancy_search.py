from uuid import UUID

from sqlalchemy import and_, select
from sqlalchemy.exc import DBAPIError, SQLAlchemyError
from sqlalchemy.orm import selectinload

from src.core.exc import InternalError
from src.domain.models import TagORM, VacanciesTagsORM, VacancyORM
from src.domain.models.vacancy import VacancyStatus
from src.infrastructure.persistence.postgresql_repo.unit_of_work import UnitOfWork
from src.utils.catch_error import catch_raise_error


class VacancySearchRepo:
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
            .where(and_(VacancyORM.status == VacancyStatus.PUBLISHED, VacancyORM.author_name.contains(author)))
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
            .where(VacancyORM.author_name.contains(author))
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

    @catch_raise_error((SQLAlchemyError, DBAPIError), InternalError, "critical", "can't find vacancies for admin by title")
    async def find_vacancies_for_admin_by_title(
        self,
        uof: UnitOfWork,
        title: str,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None:
        query = (
            select(VacancyORM)
            .where(VacancyORM.title.contains(title))
            .offset(offset)
            .limit(limit)
            .options(selectinload(VacancyORM.tags))
        )

        res = await uof.session.execute(query)
        vacancies = res.scalars().all()

        return list(vacancies) if len(vacancies) > 0 else None

    @catch_raise_error((SQLAlchemyError, DBAPIError), InternalError, "critical", "can't find vacancies for admin by title")
    async def find_only_published_vacancies_by_title(
        self,
        uof: UnitOfWork,
        title: str,
        offset: int,
        limit: int,
    ) -> list[VacancyORM] | None:
        query = (
            select(VacancyORM)
            .where(and_(VacancyORM.status == VacancyStatus.PUBLISHED, VacancyORM.title.contains(title)))
            .offset(offset)
            .limit(limit)
            .options(selectinload(VacancyORM.tags))
        )

        res = await uof.session.execute(query)
        vacancies = res.scalars().all()

        return list(vacancies) if len(vacancies) > 0 else None
