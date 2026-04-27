from sqlalchemy import select, and_
from sqlalchemy.dialects.postgresql import insert
from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker
from sqlalchemy.orm import selectinload

from src.core.exc.internal import InternalError
from src.domain.models import TagORM, VacancyORM, VacanciesTagsORM
from src.utils.catch_error import catch_rise_error


class VacancyRepo:
    def __init__(self, session_maker: async_sessionmaker[AsyncSession]) -> None:
        self.session_maker: async_sessionmaker[AsyncSession] = session_maker

    @catch_rise_error(SQLAlchemyError, InternalError, 'critical', '''can't create new vacancy in database''')
    async def create_vacancy(self, vacancy: VacancyORM) -> int:
        async with self.session_maker() as session:
            async with session.begin():
                tags = [{'tag': tag.tag} for tag in vacancy.tags]
                await session.execute(insert(TagORM).values(tags).on_conflict_do_nothing(index_elements=['tag']))

                tags = await session.execute(select(TagORM).where(TagORM.tag.in_([t['tag'] for t in tags])))

                vacancy.tags = list(tags.scalars().all())

                session.add(vacancy)
                await session.flush()

                return vacancy.vacancy_id

    @catch_rise_error(SQLAlchemyError, InternalError, 'critical', '''can't find vacancy with given vacancy_id''')
    async def find_vacancy_by_id(self, vacancy_id: int) -> VacancyORM | None:
        async with self.session_maker() as session:
            query = select(VacancyORM).where(VacancyORM.vacancy_id == vacancy_id).options(
                selectinload(VacancyORM.tags)
            )

            res = await session.execute(query)
            return res.scalar_one_or_none()

    @catch_rise_error(SQLAlchemyError, InternalError, 'critical', '''can't find vacancies with given tags''')
    async def find_vacancies_with_tags(self, tags: list[str], limit: int, offset: int) -> list[VacancyORM]:
        async with self.session_maker() as session:
            query = select(TagORM).where(TagORM.tag.in_(tags))

            res = await session.execute(query)
            tag_ids = set(row[0] for row in res)

            query = select(VacancyORM).join(
                VacanciesTagsORM, VacanciesTagsORM.vacancy_id == VacancyORM.vacancy_id
            ).where(VacanciesTagsORM.tag_id.in_(tag_ids)).offset(offset).limit(limit)

            res = await session.execute(query)
            return list(res.scalars())
