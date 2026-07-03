from typing import Protocol


class IValidationServiceClient(Protocol):
    async def is_metro_valid(self, city: str, metro: str) -> bool: ...
