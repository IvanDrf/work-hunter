import logging

from redis.asyncio import Redis
from redis.exceptions import ConnectionError, RedisError

from src.core.exc import InternalError


class CityRedisRepo:
    def __init__(self, redis: Redis) -> None:
        self.redis: Redis = redis

        self.logger = logging.getLogger("CityRedisRepo")

    async def is_city_exists(self, city: str) -> bool:
        try:
            return await self.redis.hexists(name="cities", key=city)

        except (RedisError, ConnectionError) as e:
            self.logger.critical(f"can't check is {city=}, details={e}")

            raise InternalError(f"can't check is {city=} exists")

    async def close(self) -> None:
        await self.redis.aclose()
