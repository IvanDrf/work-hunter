from src.core.config.app import AppConfig
from src.core.config.broker import RabbitMQConfig
from src.core.config.database import PostgreSQLConfig, RedisConfig


class Config(AppConfig, PostgreSQLConfig, RabbitMQConfig, RedisConfig):
    pass
