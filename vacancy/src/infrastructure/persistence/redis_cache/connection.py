import logging

from redis.asyncio import Redis
from redis.exceptions import ConnectionError

from src.core.config.cache import RedisConfig


async def connect(config: RedisConfig) -> Redis:
    logger = logging.getLogger("PostgreSQL_connect")

    client = Redis(host=config.redis_host, port=config.redis_port, db=config.redis_db)

    try:
        await client.ping()  # type: ignore
    except ConnectionError as e:
        logger.error(f"can't connect to redis, error={e}")
        raise

    return client
