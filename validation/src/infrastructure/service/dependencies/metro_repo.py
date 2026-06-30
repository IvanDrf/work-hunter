from typing import Protocol


class IMetroRepo(Protocol):
    async def is_metro_exists(self, city: str, metro: str) -> bool: ...
