from pkg.applications_api.applications_pb2_grpc import ApplicationServiceServicer

from src.api.grpc.handlers import ApplicationHandlers
from src.core.config import Config
from src.infrastructure.persistence.postgresql_repo import ApplicationPostgreSQLRepo, UnitOfWork, connect_to_postgresql
from src.infrastructure.persistence.rabbitmq_broker import (
    ApplicationsRabbitMQProducer,
    connect_to_rabbitmq,
    declare_channel,
    declare_exchange,
)
from src.infrastructure.service.application import ApplicationService
from src.infrastructure.service.application.dependencies import IApplicationProducer


class Fabric:
    def __init__(self, config: Config) -> None:
        self.config: Config = config

    async def new_handlers(self) -> ApplicationServiceServicer:
        application_repo = ApplicationPostgreSQLRepo()
        engine, session_maker = await connect_to_postgresql(self.config)
        uof = UnitOfWork(engine, session_maker)

        application_producer = await self.new_application_producer()

        application_service = ApplicationService(
            uof,
            application_repo,  # type: ignore
            self.config.postgres_timeout,
            application_producer,
        )

        return ApplicationHandlers(application_service, self.config.service_timeout)

    async def new_application_producer(self) -> IApplicationProducer:
        conn = await connect_to_rabbitmq(self.config)
        chan = await declare_channel(conn)
        exchange = await declare_exchange(self.config, chan)

        return ApplicationsRabbitMQProducer(conn=conn, chan=chan, exchange=exchange, routing_key=self.config.rabbitmq_routing_key)
