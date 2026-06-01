import logging
from asyncio import wait_for
from datetime import timedelta
from typing import Final

from pydantic import ValidationError

from src.domain.schemas import VacancyResponseSchema
from src.infrastructure.service.dependencies import ICache


class BaseVacancyService:
    def __init__(self, cache: ICache, vacancy_ttl: timedelta, cache_timeout: float) -> None:
        self.cache: ICache = cache

        self.VACANCY_TTL: Final[timedelta] = vacancy_ttl
        self.CACHE_TIMEOUT: Final[float] = cache_timeout

    async def _save_vacancy_in_cache_by_id(self, vacancy: VacancyResponseSchema) -> None:
        vacancy_json = vacancy.model_dump_json()

        await self.cache.save(str(vacancy.vacancy_id), vacancy_json, self.VACANCY_TTL)

    async def _get_vacancy_from_cache_by_id(self, vacancy_id: int) -> VacancyResponseSchema | None:
        try:
            vacancy_json = await wait_for(self.cache.get(str(vacancy_id)), timeout=self.CACHE_TIMEOUT)
            return VacancyResponseSchema.model_validate_json(vacancy_json) if vacancy_json is not None else None

        except ValidationError as e:
            logging.error(f"Invalid vacancy json in cache, can't parse it, error={e}")

        except TimeoutError:
            logging.error(f"Can't get vacancy with {vacancy_id=} form cache, timeout error")
