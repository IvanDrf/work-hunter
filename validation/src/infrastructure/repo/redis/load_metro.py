from ijson import items
from redis.asyncio import Redis

from src.core.config import Config


async def load_cities_and_metro(config: Config, redis: Redis) -> None:
    with open(config.metro_json, "r", encoding="utf-8") as metro_file:
        objects = items(metro_file, "cities.item")

        for city in objects:
            await set_cities_and_metro(city, redis)


async def set_cities_and_metro(city: dict[str, list[str]], redis: Redis) -> None:
    for city_name, metro_stations in city.items():
        city_name = city_name.lower()
        for station in metro_stations:
            await redis.hset(name=city_name, key=station.lower(), value=1)
