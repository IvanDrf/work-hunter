from sqlalchemy import select, update
from sqlalchemy.dialects.postgresql import insert
from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker
from sqlalchemy.orm import selectinload

from src.core.exc.internal import InternalError
from src.domain.models import TagORM, VacanciesTagsORM, VacancyORM
from src.domain.models.vacancy import VacancyStatus
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
            query = select(VacancyORM)\
                .where(VacancyORM.vacancy_id == vacancy_id)\
                .options(selectinload(VacancyORM.tags))

            res = await session.execute(query)
            return res.scalar_one_or_none()

    @catch_rise_error(SQLAlchemyError, InternalError, 'critical', '''can't find vacancies with given tags''')
    async def find_vacancies_with_tags(self, tags: list[str], offset: int, limit: int) -> list[VacancyORM] | None:
        async with self.session_maker() as session:
            query = select(VacancyORM)\
                .join(VacanciesTagsORM, VacanciesTagsORM.vacancy_id == VacancyORM.vacancy_id)\
                .join(TagORM, TagORM.tag_id == VacanciesTagsORM.tag_id)\
                .where(TagORM.tag.in_(tags))\
                .offset(offset).limit(limit)\
                .options(selectinload(VacancyORM.tags))

            res = await session.execute(query)
            vacancies = list(res.scalars())

            return vacancies if len(vacancies) > 0 else None

    @catch_rise_error(SQLAlchemyError, InternalError, 'critical', '''can't update vacancy status''')
    async def set_vacancy_status(self, vacancy_id: int, status: VacancyStatus) -> None:
        async with self.session_maker() as session:
            query = update(VacancyORM)\
                .where(VacancyORM.vacancy_id == vacancy_id)\
                .values(status=status)\
                .returning(VacancyORM.vacancy_id)

            res = await session.execute(query)
            if res.one_or_none() is None:
                raise InternalError(
                    f'''can't update vacancy status with given {vacancy_id=}, does this vacancy exists?'''
                )

            await session.commit()
