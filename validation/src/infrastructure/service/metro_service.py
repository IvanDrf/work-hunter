import logging
from asyncio import wait_for

from src.core.exc import InternalError
from src.infrastructure.service.dependencies import IMetroRepo


class MetroService:
    def __init__(self, metro_repo: IMetroRepo, repo_timeout: float) -> None:
        self.metro_repo: IMetroRepo = metro_repo
        self.repo_timeout: float = repo_timeout

    async def is_metro_valid(self, city: str, metro: str) -> bool:
        try:
            return await wait_for(self.metro_repo.is_metro_exists(city.lower(), metro.lower()), timeout=self.repo_timeout)

        except TimeoutError as e:
            logging.critical(f"request to database is too long, {city=}, {metro=}, details={e}")
            raise InternalError(f"can't check is {metro=} exists in {city=}")

    async def close(self) -> None:
        await self.metro_repo.close()
