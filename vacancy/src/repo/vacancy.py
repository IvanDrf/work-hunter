from sqlalchemy.exc import SQLAlchemyError
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from src.core.exc.internal import InternalError
from src.domain.models.vacancy import VacancyORM
from src.utils.catch_error import catch_rise_error


class VacancyRepo:
    def __init__(self, session_maker: async_sessionmaker[AsyncSession]) -> None:
        self.session_maker: async_sessionmaker[AsyncSession] = session_maker

    @catch_rise_error(SQLAlchemyError, InternalError, 'critical', '''can't create new vacancy in database''')
    async def create_vacancy(self, vacancy: VacancyORM) -> int:
        async with self.session_maker() as session:
            session.add(vacancy)
            await session.commit()

            return vacancy.vacancy_id

    @catch_rise_error(SQLAlchemyError, InternalError, 'critical', '''can't find vacancy with given vacancy_id''')
    async def find_vacancy_by_id(self, vacancy_id: int) -> VacancyORM | None:
        async with self.session_maker() as session:
            vacancy = await session.get(VacancyORM, vacancy_id)

            return vacancy
