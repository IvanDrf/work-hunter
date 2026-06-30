import logging

from redis.asyncio import Redis
from redis.exceptions import ConnectionError, RedisError

from src.core.exc import InternalError


class MetroRedisRepo:
    def __init__(self, redis: Redis) -> None:
        self.redis: Redis = redis

    async def is_metro_exists(self, city: str, metro: str) -> bool:
        try:
            return await self.redis.hexists(city, metro)

        except (RedisError, ConnectionError) as e:
            logging.critical(f"can't check is {metro=} exists by {city=}, details={e}")

            raise InternalError(f"can't check is {metro=} exists by {city=}")

    async def close(self) -> None:
        await self.redis.aclose()
