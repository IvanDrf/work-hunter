from src.api.dependencies import IMetroService
from src.core.config import Config
from src.infrastructure.repo import MetroRedisRepo, connect_to_redis, load_cities_and_metro
from src.infrastructure.service.dependencies import IMetroRepo
from src.infrastructure.service.metro_service import MetroService


class Fabric:
    def __init__(self, config: Config) -> None:
        self.config: Config = config

    def new_service(self, repo: IMetroRepo) -> IMetroService:
        return MetroService(repo, self.config.redis_timeout)

    async def new_repo(self) -> IMetroRepo:
        redis = await connect_to_redis(self.config)
        await load_cities_and_metro(self.config, redis)

        return MetroRedisRepo(redis)
