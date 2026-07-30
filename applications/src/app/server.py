import logging
from concurrent.futures import ThreadPoolExecutor

from grpc.aio import Server as grpcServer
from grpc.aio import server
from grpc_reflection.v1alpha.reflection import SERVICE_NAME, enable_server_reflection
from pkg.applications_api.applications_pb2_grpc import ApplicationServiceServicer, add_ApplicationServiceServicer_to_server

from src.core.config import AppConfig


class Server:
    def __init__(self, config: AppConfig) -> None:
        self.config: AppConfig = config
        self.server: grpcServer | None = None

        self.shutdown_time = 1

    def register(self, handlers: ApplicationServiceServicer) -> None:
        self.handlers = handlers

        self.server = server(ThreadPoolExecutor(max_workers=self.config.workers))
        add_ApplicationServiceServicer_to_server(handlers, self.server)

        self.server.add_insecure_port(self.config.app_address)
        enable_server_reflection(SERVICE_NAME, self.server)

    async def run(self) -> None:
        if self.server is None:
            raise RuntimeError("servier is not registred")

        logging.info(f"Start applications service on {self.config.app_address}")
        await self.server.start()
        await self.server.wait_for_termination()

    async def stop(self) -> None:
        if self.server is None:
            raise RuntimeError("server is not registred")

        await self.server.stop(self.shutdown_time)
        if hasattr(self.handlers, "stop"):
            await self.handlers.stop()  # type: ignore
