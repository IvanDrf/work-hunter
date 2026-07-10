import logging
from concurrent.futures import ThreadPoolExecutor
from typing import Final

from grpc.aio import server
from grpc_reflection.v1alpha.reflection import SERVICE_NAME, enable_server_reflection
from pkg.vacancy_api.vacancy_pb2_grpc import add_VacancyServicer_to_server
from prometheus_client import start_http_server
from py_async_grpc_prometheus.prometheus_async_server_interceptor import PromAsyncServerInterceptor

from src.api.grpc.handlers import VacancyHandlers
from src.core.config import AppConfig


class Server:
    WORKERS: Final[int] = 4

    def __init__(self, config: AppConfig) -> None:
        self.address: str = config.address
        self.server_port: int = config.app_port
        self.metrics_port: int = config.metrics_port

        self.server = None
        self.metrics_server = None
        self.shutdown_time: float = config.app_shutdown_time

    def register(self, handlers: VacancyHandlers) -> None:
        self.server = server(
            ThreadPoolExecutor(max_workers=self.WORKERS),
            interceptors=(PromAsyncServerInterceptor(enable_handling_time_histogram=True),),
        )
        add_VacancyServicer_to_server(handlers, self.server)

        self.server.add_insecure_port(self.address)
        enable_server_reflection(SERVICE_NAME, self.server)

    async def run(self) -> None:
        if self.server is None:
            raise RuntimeError("server is not registred")

        logging.info(f"Starting server {self.address}")
        await self.server.start()

        logging.info(f"Starting metrics server {self.metrics_port}")
        self.metrics_server, _ = start_http_server(self.metrics_port)

        await self.server.wait_for_termination()

    async def stop(self) -> None:
        logging.info(f"Stopping metrics server {self.metrics_port}")
        if self.metrics_server is None:
            raise RuntimeError("metrics server is not registred")

        self.metrics_server.shutdown()

        if self.server is None:
            raise RuntimeError("server is not registred")

        logging.info(f"Stopping server {self.address}")
        await self.server.stop(self.shutdown_time)
