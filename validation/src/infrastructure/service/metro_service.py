from src.core.exc import ArgumentError
from src.infrastructure.service.dependencies import IMetroRepo


class MetroService:
    def __init__(self, metro_repo: IMetroRepo) -> None:
        self.metro_repo: IMetroRepo = metro_repo

    async def is_metro_valid(self, city: str, metro: str) -> bool:
        if not city or not metro:
            raise ArgumentError("city or metro is empty")

        return await self.metro_repo.is_metro_exists(city.lower(), metro.lower())
