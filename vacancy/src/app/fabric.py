from sqlalchemy.ext.asyncio import AsyncEngine, AsyncSession, async_sessionmaker

from src.api.grpc.handlers import VacancyHandlers
from src.api.rabbitmq.connection import connect_to_rabbitmq, declare_channel, declare_exchange, declare_queue
from src.api.rabbitmq.consumer import RabbitMQConsumer
from src.api.rabbitmq.dependencies import IApplicationService
from src.core.config import Config
from src.infrastructure.clients import ValidationClient
from src.infrastructure.persistence.postgresql_repo import TagRepo, UnitOfWork, VacancyRepo
from src.infrastructure.persistence.postgresql_repo import connect as connect_postgresql
from src.infrastructure.persistence.redis_cache import RedisCache
from src.infrastructure.persistence.redis_cache import connect as connect_redis
from src.infrastructure.service import ApplicationService, VacancyService
from src.infrastructure.service.dependencies import ICache, ITagRepo, IUnitOfWork, IVacancyRepo, IValidationServiceClient


class Fabric:
    def __init__(self, config: Config) -> None:
        self.config: Config = config

    async def new_handlers(self) -> VacancyHandlers:
        engine, session_maker = await connect_postgresql(self.config)

        vacancy_repo: IVacancyRepo = self.new_vacancy_repo()  # type: ignore
        tag_repo: ITagRepo = self.new_tag_repo()  # type: ignore
        cache: ICache = await self.new_cache()
        uof: IUnitOfWork = self.new_unit_of_work(engine, session_maker)

        validation_client = self.new_validation_client()

        vacancy_service = self.new_vacancy_service(vacancy_repo, tag_repo, cache, validation_client, uof)

        return VacancyHandlers(vacancy_service)

    async def new_rabbitmq_consumer(self) -> RabbitMQConsumer:
        engine, session_maker = await connect_postgresql(self.config)
        uof = self.new_unit_of_work(engine, session_maker)
        vacancy_repo = self.new_vacancy_repo()

        application_service: IApplicationService = ApplicationService(vacancy_repo, uof)  # type: ignore

        conn = await connect_to_rabbitmq(self.config)
        chan = await declare_channel(conn)

        exchange = await declare_exchange(self.config, chan)
        consumer_queue = await declare_queue(self.config, chan, exchange)

        return RabbitMQConsumer(conn, chan, consumer_queue, application_service, self.config.rabbitmq_service_timeout)

    def new_vacancy_repo(self) -> VacancyRepo:
        return VacancyRepo()

    def new_tag_repo(self) -> TagRepo:
        return TagRepo()

    async def new_cache(self) -> RedisCache:
        conn = await connect_redis(self.config)
        return RedisCache(conn)

    def new_unit_of_work(self, engine: AsyncEngine, session_maker: async_sessionmaker[AsyncSession]) -> UnitOfWork:
        return UnitOfWork(engine, session_maker)

    def new_validation_client(self) -> IValidationServiceClient:
        return ValidationClient(self.config.validation_service_url, self.config.validation_service_api_key)

    def new_vacancy_service(
        self,
        vacancy_repo: IVacancyRepo,
        tag_repo: ITagRepo,
        cache: ICache,
        validation_client: IValidationServiceClient,
        uof: IUnitOfWork,
    ) -> VacancyService:
        return VacancyService(
            vacancy_repo,
            tag_repo,
            uof,
            cache,
            validation_client,
            self.config.redis_ttl,
            self.config.redis_timeout,
        )
