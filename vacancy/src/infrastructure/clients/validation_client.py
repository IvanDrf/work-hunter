import logging
from asyncio import sleep
from functools import cached_property, wraps

from httpx import AsyncClient, ConnectError, codes

from src.core.exc import RetryError


def max_retries(*, max_retries: int, base_delay: float):
    def decorator(func):
        @wraps(func)
        async def wrapper(*args, **kwargs):
            last_exc = None

            for attempt in range(max_retries + 1):
                delay = base_delay * (attempt + 1)
                try:
                    return await func(*args, **kwargs)
                except (RetryError, ConnectError) as e:
                    logging.error(f"retrying to send request, attempt: {attempt}, details={e}")
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

    @cached_property
    def validate_metro_url(self) -> str:
        return f"{self.url}/api/metro"

    async def is_metro_valid(self, city: str, metro: str) -> bool:
        async with AsyncClient(timeout=5) as client:
            try:
                return await self._send_request_to_validation_service(client, city, metro)
            except (RetryError, ConnectError) as e:
                logging.error(f"can't validate metro after all attempts, details={e}")
                return False

    @max_retries(max_retries=3, base_delay=0.2)
    async def _send_request_to_validation_service(self, client: AsyncClient, city: str, metro: str) -> bool:
        response = await client.post(
            url=self.validate_metro_url,
            headers={"Content-type": "application/json", "api-key": self.api_key},
            json={"city": city, "metro": metro},
        )

        if response.status_code == codes.NO_CONTENT:
            return True

        if response.status_code == codes.NOT_FOUND:
            return False

        code = response.status_code
        message = response.json()
        await response.aclose()

        match code:
            case codes.UNPROCESSABLE_ENTITY:
                logging.critical(f"invalid request body in validation service, details={message}")

            case codes.FORBIDDEN:
                logging.critical(f"invalid api-key, can't validate metro, details={message}")

            case codes.UNAUTHORIZED:
                logging.critical(f"api-key wasn't send in request, details={message}")

            case codes.INTERNAL_SERVER_ERROR:
                logging.error(f"validation service error, details={message}")
                raise RetryError("validation service is not available")

            case codes.TOO_MANY_REQUESTS:
                logging.error(f"validation service error, details={message}")
                raise RetryError("too many requests to validation service")

        return False
