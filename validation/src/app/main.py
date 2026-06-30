from contextlib import asynccontextmanager

from fastapi import FastAPI
from uvicorn import run

from src.api.dependencies import init_metro_service
from src.api.health import health_router
from src.api.metro import metro_router
from src.app.fabric import Fabric
from src.core.config import Config
from src.core.logger import setup_logger


@asynccontextmanager
async def lifespan(_: FastAPI):
    fabric = Fabric(config)
    setup_logger(config.logger_level)

    metro_service = fabric.new_service(await fabric.new_repo())
    init_metro_service(metro_service)

    yield

    await metro_service.close()


config = Config()

app = FastAPI(lifespan=lifespan)
app.include_router(health_router)
app.include_router(metro_router)


if __name__ == "__main__":
    run(app=app, host=config.app_host, port=config.app_port)
