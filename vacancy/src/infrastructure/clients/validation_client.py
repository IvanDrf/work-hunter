import logging
from asyncio import sleep
from functools import cached_property, wraps

from httpx import AsyncClient, ConnectError, codes

from src.core.exc import RetryError


def max_retries(*, max_retries: int, base_delay: float):
    def decorator(func):
        logger = logging.getLogger(f"max_retries: {func.__name__}")

        @wraps(func)
        async def wrapper(*args, **kwargs):
            last_exc = None

            for attempt in range(max_retries + 1):
                delay = base_delay * (attempt + 1)
                try:
                    return await func(*args, **kwargs)
                except (RetryError, ConnectError) as e:
                    logger.error(f"retrying to send request, attempt: {attempt}, details={e}")
                    await sleep(delay)

                    last_exc = e

            if isinstance(last_exc, Exception):
                raise last_exc

        return wrapper

    return decorator


class ValidationClient:
    def __init__(self, url: str, api_key: str) -> None:
        self.url: str = url
        self.api_key: str = api_key

        self.logger = logging.getLogger("ValidationClient")

    @cached_property
    def validate_metro_url(self) -> str:
        return f"{self.url}/api/metro"

    @cached_property
    def validate_city_url(self) -> str:
        return f"{self.url}/api/city"

    async def is_metro_valid(self, city: str, metro: str) -> bool:
        async with AsyncClient(timeout=5) as client:
            try:
                return await self._send_request_to_validation_service(
                    client=client,
                    url=self.validate_metro_url,
                    json={"city": city, "metro": metro},
                )

            except (RetryError, ConnectError) as e:
                self.logger.error(f"can't validate metro after all attempts, details={e}")
                return False

    async def is_city_valid(self, city: str) -> bool:
        async with AsyncClient(timeout=5) as client:
            try:
                return await self._send_request_to_validation_service(
                    client=client,
                    url=self.validate_city_url,
                    json={"city": city},
                )

            except (RetryError, ConnectError) as e:
                self.logger.error(f"can't validate city after all attempts, details={e}")
                return False

    @max_retries(max_retries=3, base_delay=0.2)
    async def _send_request_to_validation_service(self, client: AsyncClient, url: str, json: dict) -> bool:
        response = await client.post(
            url=url,
            headers={"Content-type": "application/json", "api-key": self.api_key},
            json=json,
        )

        if response.status_code == codes.NO_CONTENT:
            await response.aclose()
            return True

        if response.status_code == codes.NOT_FOUND:
            await response.aclose()
            return False

        code = response.status_code
        message = response.json()
        await response.aclose()

        match code:
            case codes.UNPROCESSABLE_ENTITY:
                self.logger.critical(f"invalid request body in validation service, details={message}, {json=}")

            case codes.FORBIDDEN:
                self.logger.critical(f"invalid api-key, can't validate, details={message}")

            case codes.UNAUTHORIZED:
                self.logger.critical(f"api-key wasn't send in request, details={message}")

            case codes.INTERNAL_SERVER_ERROR:
                self.logger.error(f"validation service error, details={message}")
                raise RetryError("validation service is not available")

            case codes.TOO_MANY_REQUESTS:
                self.logger.error(f"validation service error, details={message}")
                raise RetryError("too many requests to validation service")

        return False
