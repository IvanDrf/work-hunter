import logging
from concurrent.futures import ThreadPoolExecutor
from typing import Final

from grpc.aio import server
from grpc_reflection.v1alpha.reflection import SERVICE_NAME, enable_server_reflection
from pkg.vacancy_api.vacancy_pb2_grpc import add_VacancyServicer_to_server

from src.api.grpc.handlers import VacancyHandlers
from src.core.config import AppConfig


class Server:
    WORKERS: Final[int] = 4

    def __init__(self, config: AppConfig) -> None:
        self.address: str = config.address

        self.server = None
        self.shutdown_time: float = config.app_shutdown_time

    def register(self, handlers: VacancyHandlers) -> None:
        self.server = server(ThreadPoolExecutor(max_workers=self.WORKERS))
        add_VacancyServicer_to_server(handlers, self.server)

        self.server.add_insecure_port(self.address)
        enable_server_reflection(SERVICE_NAME, self.server)

    async def run(self) -> None:
        if self.server is None:
            raise RuntimeError("server is not registred")

        logging.info(f"Starting server {self.address}")

        await self.server.start()
        await self.server.wait_for_termination()

    async def stop(self) -> None:
        if self.server is None:
            raise RuntimeError("server is not registred")

        logging.info(f"Stopping server {self.address}")
        await self.server.stop(self.shutdown_time)
