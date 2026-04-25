from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine
from sqlalchemy import text

from src.core.config.database import PostgreSQLConfig


async def connect(config: PostgreSQLConfig) -> async_sessionmaker[AsyncSession]:
    engine = create_async_engine(config.dsn, pool_pre_ping=True)

    async with engine.begin() as conn:
        await conn.execute(text('SELECT 1'))

    return async_sessionmaker(engine)
