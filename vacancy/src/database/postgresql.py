from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine

from src.core.config.database import PostgreSQLConfig


def connect(config: PostgreSQLConfig) -> async_sessionmaker[AsyncSession]:
    engine = create_async_engine(config.dsn)

    return async_sessionmaker(engine)
