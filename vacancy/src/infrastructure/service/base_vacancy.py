import logging
from asyncio import wait_for
from datetime import timedelta
from typing import Final

from pydantic import TypeAdapter, ValidationError

from src.domain.schemas import UserInfo, VacancyResponseSchema
from src.domain.types import OrderBy
from src.infrastructure.service.dependencies import ICache

Vacancies = TypeAdapter(list[VacancyResponseSchema])


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

    async def _save_vacancies_by_author(
        self,
        vacancies: list[VacancyResponseSchema],
        author: str,
        order_by: OrderBy,
        user_info: UserInfo | None,
        offset: int,
        limit: int,
    ) -> None:
        vacancies_json = Vacancies.dump_json(vacancies).decode()
        key = generate_cache_key(author, order_by, user_info, offset, limit)

        await self.cache.save(key, vacancies_json, self.VACANCY_TTL)

    async def _get_vacancies_by_author(
        self,
        author: str,
        order_by: OrderBy,
        user_info: UserInfo | None,
        offset: int,
        limit: int,
    ) -> list[VacancyResponseSchema] | None:
        try:
            key = generate_cache_key(author, order_by, user_info, offset, limit)
            vacancies_json = await wait_for(self.cache.get(key), timeout=self.CACHE_TIMEOUT)

            return Vacancies.validate_json(vacancies_json) if vacancies_json is not None else None

        except ValidationError as e:
            logging.error(f"Invalid vacancies json in cache, can't parse it, error={e}")

        except TimeoutError:
            logging.error(f"Can't get vacancy by author with args {author=}, {offset=}, {limit=} form cache, timeout error")


def generate_cache_key(author: str, order_by: OrderBy, user_info: UserInfo | None, offset: int, limit: int) -> str:
    return f"{author}:{order_by}:{user_info.user_role if user_info else None}:{offset}:{limit}"
