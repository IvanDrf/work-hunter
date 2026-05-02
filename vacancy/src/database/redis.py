import logging

from redis.asyncio import Redis
from redis.exceptions import ConnectionError
from src.core.config.cache import RedisConfig


async def connect(config: RedisConfig) -> Redis | None:
    client = Redis(
        host=config.host,
        port=config.port,
        db=config.db
    )

    try:
        await client.ping()  # type: ignore
    except ConnectionError as e:
        logging.error(f'''can't connect to redis, details={e}''')
        return None

    return client
