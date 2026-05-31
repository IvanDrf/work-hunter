from pydantic import Field

from src.core.config.base import BaseConfig


class AppConfig(BaseConfig):
    app_host: str = Field(default="localhost", validation_alias="APP_HOST")
    app_port: int = Field(default=50052, validation_alias="APP_PORT")

    logger_level: str = Field(default="INFO", validation_alias="LOGGER_LEVEL")

    @property
    def address(self) -> str:
        return f"{self.app_host}:{self.app_port}"
