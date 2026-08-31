from typing import Self

from sqlalchemy.ext.asyncio import AsyncEngine, AsyncSession, async_sessionmaker


class UnitOfWork:
    def __init__(self, engine: AsyncEngine, session_maker: async_sessionmaker[AsyncSession]) -> None:
        self.engine: AsyncEngine = engine
        self.session_maker: async_sessionmaker[AsyncSession] = session_maker

    async def __aenter__(self) -> Self:
        self.session: AsyncSession = self.session_maker()
        return self

    async def __aexit__(self, exc_type, exc, tb):
        if exc_type:
            await self.rollback()
        else:
            await self.commit()

    async def commit(self) -> None:
        await self.session.commit()

    async def rollback(self) -> None:
        await self.session.rollback()

    async def flush(self) -> None:
        await self.session.flush()

    async def stop(self) -> None:
        if hasattr(self, "session") and self.session is not None:
            await self.session.close_all()

        await self.engine.dispose()
