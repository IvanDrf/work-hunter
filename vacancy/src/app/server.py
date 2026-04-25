from concurrent.futures import ThreadPoolExecutor
from typing import Final

from grpc.aio import server
from grpc_reflection.v1alpha.reflection import SERVICE_NAME, enable_server_reflection

from src.core.config.app import AppConfig
from src.core.exc.internal import InternalError


class Server:
    WORKERS: Final[int] = 4

    def __init__(self, config: AppConfig) -> None:
        self.host: str = config.host
        self.port: int = config.port
        self.server = None

    def register(self, handlers) -> None:
        self.server = server(
            ThreadPoolExecutor(max_workers=self.WORKERS), handlers
        )
        self.server.add_insecure_port(f'{self.host}:{self.port}')

        enable_server_reflection(SERVICE_NAME, self.server)

    async def run(self) -> None:
        if self.server is None:
            raise InternalError('server is not registred')

        await self.server.start()
        await self.server.wait_for_termination()
