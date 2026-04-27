from dataclasses import dataclass


@dataclass(frozen=True, slots=True, kw_only=True)
class AppConfig:
    host: str
    port: int

    logger_level:  str

    @property
    def address(self) -> str:
        return f'{self.host}:{self.port}'
