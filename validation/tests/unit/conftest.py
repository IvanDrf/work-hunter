from fakeredis.aioredis import FakeRedis
from pytest_asyncio import fixture as async_fixture

from src.core.config import Config
from src.infrastructure.repo import MetroRedisRepo, load_cities_and_metro


@async_fixture(scope="session")
async def metro_repo() -> MetroRedisRepo:
    redis = FakeRedis()

    await load_cities_and_metro(Config(), redis)

    return MetroRedisRepo(redis)
