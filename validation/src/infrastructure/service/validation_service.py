import logging
from asyncio import wait_for

from src.core.exc import InternalError
from src.infrastructure.service.dependencies import ICityRepo, IMetroRepo


class ValidationService:
    def __init__(self, metro_repo: IMetroRepo, city_repo: ICityRepo, repo_timeout: float) -> None:
        self.metro_repo: IMetroRepo = metro_repo
        self.city_repo: ICityRepo = city_repo

        self.repo_timeout: float = repo_timeout

    async def is_metro_valid(self, city: str, metro: str) -> bool:
        try:
            return await wait_for(self.metro_repo.is_metro_exists(city.lower(), metro.lower()), timeout=self.repo_timeout)

        except TimeoutError as e:
            logging.critical(f"request to database is too long, {city=}, {metro=}, details={e}")
            raise InternalError(f"can't check is {metro=} exists in {city=}")

    async def is_city_valid(self, city: str) -> bool:
        try:
            return await wait_for(self.city_repo.is_city_exists(city.lower()), timeout=self.repo_timeout)
        except TimeoutError as e:
            logging.critical(f"request to database is too long, {city=}, details={e}")
            raise InternalError(f"can't check is {city=} exists")

    async def close(self) -> None:
        await self.metro_repo.close()
        await self.city_repo.close()
