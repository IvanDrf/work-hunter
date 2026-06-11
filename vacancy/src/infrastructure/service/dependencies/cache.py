from datetime import timedelta
from typing import Protocol


class ICache(Protocol):
    async def save(self, key: str, content: str, ttl: timedelta) -> None: ...
    async def get(self, key: str) -> str | None: ...
