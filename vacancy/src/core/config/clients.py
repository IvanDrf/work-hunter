from pydantic import Field

from src.core.config.base import BaseConfig


class ValidationServiceConfig(BaseConfig):
    validation_service_url: str = Field(default="", validation_alias="VALIDATION_SERVICE_URL")
    validation_service_api_key: str = Field(default="", validation_alias="VALIDATION_SERVICE_API_KEY")
