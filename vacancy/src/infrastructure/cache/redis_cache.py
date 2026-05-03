from datetime import timedelta
import logging

from redis.asyncio import Redis
from redis.exceptions import RedisError


class RedisCache:
    def __init__(self, client: Redis) -> None:
        self.client: Redis = client

    async def save(self, key: str, content: str, ttl: timedelta) -> None:
        try:
            await self.client.set(key, content, ex=ttl)
        except RedisError as e:
            logging.error(f'''can't save in cache by {key=}, details={e} ''')

    async def get(self, key: str) -> str | None:
        try:
            return await self.client.get(key)
        except RedisError as e:
            logging.error(f'''can't get content from cache by {key=}, details={e} ''')
