from src.core.config.app import AppConfig
from src.core.config.broker import RabbitMQConfig
from src.core.config.database import PostgreSQLConfig


class Config(AppConfig, PostgreSQLConfig, RabbitMQConfig):
    pass
