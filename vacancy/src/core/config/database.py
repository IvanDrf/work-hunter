from dataclasses import dataclass


@dataclass(frozen=True, slots=True, kw_only=True)
class PostgreSQLConfig:
    host: str
    port: str

    username: str
    password: str

    db_name: str

    @property
    def address(self) -> str:
        return f'{self.host}:{self.port}'

    @property
    def dsn(self) -> str:
        return f'postgresql+asyncpg://{self.username}:{self.password}@{self.address}/{self.db_name}'
