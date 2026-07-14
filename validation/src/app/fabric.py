from asyncio import gather

from redis.asyncio import Redis

from src.api.dependencies import IValidationService
from src.core.config import Config
from src.infrastructure.repo import CityRedisRepo, MetroRedisRepo, connect_to_redis, load_cities, load_cities_and_metro
from src.infrastructure.service import ValidationService
from src.infrastructure.service.dependencies import ICityRepo, IMetroRepo


class Fabric:
    def __init__(self, config: Config) -> None:
        self.config: Config = config

    async def new_service(self) -> IValidationService:
        redis = await connect_to_redis(self.config)

        repos = await gather(*[self.new_metro_repo(redis), self.new_city_repo(redis)])

        return self._new_service(repos[0], repos[1])

    def _new_service(self, metro_repo: IMetroRepo, city_repo: ICityRepo) -> IValidationService:
        return ValidationService(metro_repo, city_repo, self.config.redis_timeout)

    async def new_metro_repo(self, redis: Redis) -> IMetroRepo:
        await load_cities_and_metro(self.config, redis)

        return MetroRedisRepo(redis)

    async def new_city_repo(self, redis: Redis) -> ICityRepo:
        await load_cities(self.config, redis)

        return CityRedisRepo(redis)
