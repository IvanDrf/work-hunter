from pydantic import Field

from src.core.config.base import BaseConfig


class PostgreSQLConfig(BaseConfig):
    postgres_host: str = Field(default="localhost", validation_alias="POSTGRES_HOST")
    postgres_port: int = Field(default=5432, validation_alias="POSTGRES_PORT")

    postgres_user: str = Field(default="user", validation_alias="POSTGRES_USER")
    postgres_password: str = Field(default="password", validation_alias="POSTGRES_PASSWORD")

    postgres_database: str = Field(default="database", validation_alias="POSTGRES_DATABASE")
    postgres_timeout: int = Field(default=2, validation_alias="POSTGRES_TIMEOUT")

    @property
    def postgres_address(self) -> str:
        return f"{self.postgres_host}:{self.postgres_port}"

    @property
    def postgres_dsn(self) -> str:
        return (
            f"postgresql+asyncpg://{self.postgres_user}:{self.postgres_password}@{self.postgres_address}/{self.postgres_database}"
        )


class RedisConfig(BaseConfig):
    redis_host: str = Field(default="localhost", validation_alias="REDIS_HOST")
    redis_port: int = Field(default=6379, validation_alias="REDIS_PORT")

    redis_database: int = Field(default=1, validation_alias="REDIS_DATABASE")
