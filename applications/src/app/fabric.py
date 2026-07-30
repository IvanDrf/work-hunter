from pkg.applications_api.applications_pb2_grpc import ApplicationServiceServicer
from sqlalchemy.ext.asyncio import AsyncEngine, AsyncSession, async_sessionmaker

from src.api.grpc.dependencies import IApplicationService
from src.api.grpc.handlers import ApplicationHandlers
from src.core.config import Config
from src.infrastructure.persistence.postgresql_repo import ApplicationPostgreSQLRepo, UnitOfWork, connect_to_postgresql
from src.infrastructure.service.application import ApplicationService
from src.infrastructure.service.application.dependencies import IApplicationRepo, IUnitOfWork


class Fabric:
    def __init__(self, config: Config) -> None:
        self.config: Config = config

    async def new_handlers(self) -> ApplicationServiceServicer:
        application_repo = self.new_postgresql_repo()
        engine, session_maker = await connect_to_postgresql(self.config)
        uof = self.new_unit_of_work(engine, session_maker)

        application_service = self.new_application_service(uof, application_repo)

        return ApplicationHandlers(application_service, self.config.service_timeout)

    def new_application_service(self, uof: IUnitOfWork, repo: IApplicationRepo) -> IApplicationService:
        return ApplicationService(uof, repo, self.config.postgres_timeout)

    def new_unit_of_work(self, engine: AsyncEngine, session_maker: async_sessionmaker[AsyncSession]) -> IUnitOfWork:
        return UnitOfWork(engine, session_maker)

    def new_postgresql_repo(self) -> IApplicationRepo:
        return ApplicationPostgreSQLRepo()  # type: ignore cuz double protocols: IRepo and IUnitOfWork
