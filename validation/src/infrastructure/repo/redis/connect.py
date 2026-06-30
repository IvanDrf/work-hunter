from redis.asyncio import Redis

from src.core.config import RedisConfig


async def connect_to_redis(config: RedisConfig) -> Redis:
    redis = Redis(host=config.redis_host, port=config.redis_port, db=config.redis_database)

    if not await redis.ping():
        raise TimeoutError(f"can't connect to redis, {config=}")

    return redis
