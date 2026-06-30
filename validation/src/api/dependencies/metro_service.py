from typing import Protocol


class IMetroService(Protocol):
    async def is_metro_valid(self, city: str, metro: str) -> bool: ...


def get_metro_service() -> IMetroService: ...
