import logging

from sqlalchemy.ext.asyncio import async_sessionmaker, AsyncSession
from sqlalchemy.exc import SQLAlchemyError

from src.domain.models.vacancy import VacancyORM
from src.core.exc.internal import InternalError


class VacancyRepo:
    def __init__(self, session_maker: async_sessionmaker[AsyncSession]) -> None:
        self.session_maker: async_sessionmaker[AsyncSession] = session_maker

    async def create_vacancy(self, vacancy: VacancyORM) -> int:
        try:
            return await self._create_vacancy(vacancy)
        except SQLAlchemyError as e:
            logging.critical(f'create_vacancy repo: {e}')
            raise InternalError('''can't create new vacancy in database''')

    async def _create_vacancy(self, vacancy: VacancyORM) -> int:
        async with self.session_maker() as session:
            session.add(vacancy)
            await session.commit()

            return vacancy.vacancy_id
