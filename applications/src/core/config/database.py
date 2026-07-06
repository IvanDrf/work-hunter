from pydantic import Field
from src.core.config.base import BaseConfig


class PostgreSQLConfig(BaseConfig):
    postgres_host: str = Field(default="localhost", validation_alias="POSTGRES_HOST")
    postgres_port: int = Field(default=5432, validation_alias="POSTGRES_PORT")

    postgres_user: str = Field(default="user", validation_alias="POSTGRES_USER")
    postgres_password: str = Field(default="password", validation_alias="POSTGRES_PASSWORD")

    postgres_database: str = Field(default="database", validation_alias="POSTGRES_DATABASE")
    postgres_timeout: int = Field(default=2, validation_alias="POSTGRES_TIMEOUT")
