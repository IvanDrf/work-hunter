from src.api.rabbitmq.consumer import RabbitMQConsumer
from src.app.fabric import Fabric
from src.app.server import Server
from src.core.config import Config


class ServerApp:
    def __init__(self, config: Config) -> None:
        self.server: Server = Server(config)

        self.fabric: Fabric = Fabric(config)

    async def init(self) -> None:
        self.handlers = await self.fabric.new_handlers()

        self.server.register(self.handlers)

    async def run(self) -> None:
        await self.server.run()

    async def stop(self) -> None:
        await self.server.stop()
        await self.handlers.stop()


class ConsumerApp:
    def __init__(self, config: Config) -> None:
        self.config: Config = config
        self.fabric: Fabric = Fabric(config)

        self.rabbitmq_consumer: RabbitMQConsumer | None = None

    async def init(self) -> None:
        self.rabbitmq_consumer = await self.fabric.new_rabbitmq_consumer()

    async def run(self) -> None:
        if self.rabbitmq_consumer is None:
            raise RuntimeError("RabbitMQ is not initialized")

        await self.rabbitmq_consumer.start_consuming()
