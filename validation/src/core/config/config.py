from src.core.config.app import AppConfig
from src.core.config.database import RedisConfig


class Config(AppConfig, RedisConfig):
    pass
