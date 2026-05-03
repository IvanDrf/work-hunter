from typing import Protocol
from datetime import timedelta


class ICache(Protocol):
    async def save(self, key: str, content: str, ttl: timedelta) -> None:
        ...

    async def get(self, key: str) -> str | None:
        ...
