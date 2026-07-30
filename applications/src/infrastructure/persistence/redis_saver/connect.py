from redis.asyncio.client import Redis

from src.core.config import RedisConfig


def connect_to_redis(config: RedisConfig) -> Redis:
    return Redis(host=config.redis_host, port=config.redis_port, db=config.redis_database)
