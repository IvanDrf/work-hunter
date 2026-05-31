from pydantic import Field

from src.core.config.base import BaseConfig


class RedisConfig(BaseConfig):
    redis_host: str = Field(default="localhost", validation_alias="REDIS_HOST")
    redis_port: int = Field(default=6379, validation_alias="REDIS_PORT")
    redis_db: int = Field(default=0, validation_alias="REDIS_DB")
    redis_ttl: int = Field(default=5, validation_alias="REDIS_CACHE_TTL")
