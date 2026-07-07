from pydantic import Field
from src.core.config.base import BaseConfig


class AppConfig(BaseConfig):
    app_host: str = Field(default="localhost", validation_alias="APP_HOST")
    app_port: int = Field(default=50053, validation_alias="APP_PORT")

    service_timeout: float = Field(default=2, validation_alias="SERVICE_TIMEOUT")
    workers: int = Field(default=4, gt=0, validation_alias="WORKERS")

    @property
    def app_address(self) -> str:
        return f"{self.app_host}:{self.app_port}"
