from sqlalchemy import text
from sqlalchemy.ext.asyncio import AsyncEngine, AsyncSession, async_sessionmaker, create_async_engine
from src.core.config import PostgreSQLConfig


async def connect_to_postgresql(config: PostgreSQLConfig) -> tuple[AsyncEngine, async_sessionmaker[AsyncSession]]:
    engine = create_async_engine(config.postgres_dsn)

    await ping_database(engine)

    return engine, async_sessionmaker(engine)


async def ping_database(engine: AsyncEngine) -> None:
    async with engine.begin() as conn:
        await conn.execute(text("SELECT 1"))
