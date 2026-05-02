from dataclasses import dataclass


@dataclass(slots=True, frozen=True, kw_only=True)
class RedisConfig:
    host: str
    port: int
    db: int
