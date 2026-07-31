from typing import Protocol


class IValidationService(Protocol):
    async def is_metro_valid(self, city: str, metro: str) -> bool: ...
    async def is_city_valid(self, city: str) -> bool: ...

    async def close(self) -> None: ...


validation_service: IValidationService | None = None


def init_validation_service(service: IValidationService) -> None:
    global validation_service
    validation_service = service


def get_validation_service() -> IValidationService:
    if validation_service is None:
        raise RuntimeError("metro service is None")

    return validation_service
