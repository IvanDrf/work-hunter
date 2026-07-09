from contextlib import asynccontextmanager

from fastapi import FastAPI
from prometheus_fastapi_instrumentator.instrumentation import PrometheusFastApiInstrumentator
from uvicorn import run

from src.api.dependencies import init_validation_service
from src.api.health import health_router
from src.api.validation import validation_router
from src.app.fabric import Fabric
from src.core.config import Config
from src.core.logger import setup_logger


@asynccontextmanager
async def lifespan(_: FastAPI):
    fabric = Fabric(config)
    setup_logger(config.logger_level)

    validation_service = await fabric.new_service()
    init_validation_service(validation_service)

    yield

    await validation_service.close()


config = Config()

app = FastAPI(lifespan=lifespan)
app.include_router(health_router)
app.include_router(validation_router)

PrometheusFastApiInstrumentator().instrument(app).expose(app)

if __name__ == "__main__":
    run(app=app, host=config.app_host, port=config.app_port)
