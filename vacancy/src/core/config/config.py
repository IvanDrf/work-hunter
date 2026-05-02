from dataclasses import dataclass
from pathlib import Path

from yaml import safe_load

from src.core.config.app import AppConfig
from src.core.config.database import PostgreSQLConfig
from src.core.config.cache import RedisConfig


@dataclass(frozen=True, slots=True)
class Config:
    app: AppConfig
    database: PostgreSQLConfig
    cache: RedisConfig

    @classmethod
    def load_from_yaml(cls, path: str) -> 'Config':
        p = Path(path)

        if p.is_dir() or not p.exists():
            raise FileNotFoundError(
                f'''can't file config file with given path: {path}'''
            )

        with open(p, 'r') as config_file:
            content: dict = safe_load(config_file)

            return cls(
                app=AppConfig(**content['app']),
                database=PostgreSQLConfig(**content['database']),
                cache=RedisConfig(**content['cache'])
            )
