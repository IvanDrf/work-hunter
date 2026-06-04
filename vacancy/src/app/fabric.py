from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker

from src.api.handlers import VacancyHandlers
from src.core.config import Config
from src.infrastructure.persistence.postgresql_repo import TagRepo, UnitOfWork, VacancyRepo
from src.infrastructure.persistence.postgresql_repo import connect as connect_postgresql
from src.infrastructure.persistence.redis_cache import RedisCache
from src.infrastructure.persistence.redis_cache import connect as connect_redis
from src.infrastructure.service.vacancy_service import VacancyService
from src.infrastructure.service.dependencies import ICache, ITagRepo, IUnitOfWork, IVacancyRepo


class Fabric:
    def __init__(self, config: Config) -> None:
        self.config: Config = config

    async def new_handlers(self) -> VacancyHandlers:
        session_maker = await connect_postgresql(self.config)

        vacancy_repo: IVacancyRepo = self.new_vacancy_repo()
        tag_repo: ITagRepo = self.new_tag_repo()
        cache: ICache = await self.new_cache()

        uof: IUnitOfWork = self.new_unit_of_work(session_maker)

        vacancy_service = self.new_vacancy_service(vacancy_repo, tag_repo, cache, uof)

        return VacancyHandlers(vacancy_service)

    def new_vacancy_repo(self) -> VacancyRepo:
        return VacancyRepo()

    def new_tag_repo(self) -> TagRepo:
        return TagRepo()

    async def new_cache(self) -> RedisCache:
        conn = await connect_redis(self.config)
        return RedisCache(conn)

    def new_unit_of_work(self, session_maker: async_sessionmaker[AsyncSession]) -> UnitOfWork:
        return UnitOfWork(session_maker)

    def new_vacancy_service(
        self,
        vacancy_repo: IVacancyRepo,
        tag_repo: ITagRepo,
        cache: ICache,
        uof: IUnitOfWork,
    ) -> VacancyService:
        return VacancyService(vacancy_repo, tag_repo, uof, cache, self.config.redis_ttl, self.config.redis_timeout)
