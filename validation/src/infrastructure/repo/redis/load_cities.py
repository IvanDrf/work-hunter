from ijson import items
from redis.asyncio import Redis

from src.core.config import Config


async def load_cities(config: Config, redis: Redis) -> None:
    with open(config.cities_json, "r", encoding="utf-8") as json_file:
        cities = items(json_file, "cities.item")

        for city in cities:
            await redis.hset(name="cities", key=city.lower(), value=1)
