from pydantic import Field, field_validator

from src.core.config.base import BaseConfig


class AppConfig(BaseConfig):
    app_host: str = Field(default="localhost", validation_alias="APP_HOST")
    app_port: int = Field(default=8080, validation_alias="APP_PORT")

    metro_json: str = Field(default="metro.json", validation_alias="METRO_JSON")

    api_key: str = Field(default="", validation_alias="API_KEY")

    logger_level: str = Field(default="DEBUG", validation_alias="LOGGER_LEVEL")

    @field_validator("api_key")
    def validate_api_key(cls, value: str) -> str:
        if len(value) == 0:
            raise ValueError("api-key must be non empty")

        return value
