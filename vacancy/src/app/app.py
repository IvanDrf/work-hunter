from src.app.fabric import Fabric
from src.app.server import Server
from src.core.config.config import Config


class App:
    def __init__(self, config: Config) -> None:
        self.config: Config = config
        self.server: Server = Server(config.app)

        self.fabric: Fabric = Fabric(config)

    async def init(self) -> None:
        handlers = await self.fabric.new_handlers()

        self.server.register(handlers)

    async def run(self) -> None:
        await self.server.run()
